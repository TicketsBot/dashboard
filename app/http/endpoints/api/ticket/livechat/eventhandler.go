package livechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/internal/api"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
)

func (c *Client) HandleEvent(event Event) error {
	switch event.Type {
	case EventTypeAuth:
		var data AuthData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			c.Write(NewErrorMessage("Malformed event payload"))
			_ = c.Ws.Close()
			c.Flush()
			return err
		}

		if err := c.handleAuthEvent(data); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) handleAuthEvent(data AuthData) error {
	if c.Authenticated {
		return api.NewErrorWithMessage(http.StatusBadRequest, errors.New("Already authenticated"), "Already authenticated")
	}

	token, err := jwt.Parse(data.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Conf.Server.Secret), nil
	})
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	userIdStr, ok := claims["userid"].(string)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	// Get the ticket
	ticket, err := dbclient.Client.Tickets.Get(c.TicketId, c.GuildId)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving ticket data")
	}

	if ticket.Id == 0 || ticket.GuildId == 0 || ticket.GuildId != c.GuildId {
		return api.NewErrorWithMessage(http.StatusNotFound, err, "Ticket not found")
	}

	// Verify the user has permissions to be here
	hasPermission, requestErr := utils.HasPermissionToViewTicket(c.GuildId, userId, ticket)
	if requestErr != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving permission data")
	}

	if !hasPermission {
		return api.NewErrorWithMessage(http.StatusForbidden, err, "You do not have permission to view this ticket")
	}

	// Check premium
	botContext, err := botcontext.ContextForGuild(c.GuildId)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving bot context")
	}

	// Verify the guild is premium
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(c.GuildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving premium tier")
	}

	if premiumTier == premium.None {
		return api.NewErrorWithMessage(http.StatusPaymentRequired, err, "Live-chat requires premium to use")
	}

	c.Authenticated = true

	c.Write(Event{
		Type: EventTypeAuthenticated,
	})

	return nil
}

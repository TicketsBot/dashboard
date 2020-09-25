package api

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

type sendMessageBody struct {
	Message string `json:"message"`
}

func SendMessage(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	// Get ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid ticket ID",
		})
		return
	}

	var body sendMessageBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Message is missing",
		})
		return
	}

	// Verify guild is premium
	premiumTier := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if premiumTier == premium.None {
		ctx.AbortWithStatusJSON(402, gin.H{
			"success": false,
			"error":   "Guild is not premium",
		})
		return
	}

	// Get ticket
	ticket, err := database.Client.Tickets.Get(ticketId, guildId)

	// Verify the ticket exists
	if ticket.UserId == 0 {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": false,
			"error":   "Ticket not found",
		})
		return
	}

	// Verify the user has permission to send to this guild
	if ticket.GuildId != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error":   "Guild ID doesn't match",
		})
		return
	}

	user, _ := cache.Instance.GetUser(userId)

	if len(body.Message) > 2000 {
		body.Message = body.Message[0:1999]
	}

	// Preferably send via a webhook
	webhook, err := database.Client.Webhooks.Get(guildId, ticketId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if webhook.Id != 0 {
		// TODO: Ratelimit
		_, err = rest.ExecuteWebhook(webhook.Token, nil, webhook.Id, true, rest.WebhookBody{
			Content:   body.Message,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl(256),
		})

		if err != nil {
			// We can delete the webhook in this case
			var unwrapped request.RestError
			if errors.As(err, &unwrapped); unwrapped.ErrorCode == 403 || unwrapped.ErrorCode == 404 {
				go database.Client.Webhooks.Delete(guildId, ticketId)
			}
		} else {
			ctx.JSON(200, gin.H{
				"success": true,
			})
			return
		}
	}

	body.Message = fmt.Sprintf("**%s**: %s", user.Username, body.Message)
	if len(body.Message) > 2000 {
		body.Message = body.Message[0:1999]
	}

	if ticket.ChannelId == nil {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": false,
			"error":   "Ticket channel ID is nil",
		})
		return
	}

	if _, err = rest.CreateMessage(botContext.Token, botContext.RateLimiter, *ticket.ChannelId, rest.CreateMessageData{Content: body.Message}); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

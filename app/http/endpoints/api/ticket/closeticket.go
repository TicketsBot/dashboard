package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/closerelay"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type closeBody struct {
	Reason string `json:"reason"`
}

func CloseTicket(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)
	guildId := c.Keys["guildid"].(uint64)

	ticketId, err := strconv.Atoi(c.Param("ticketId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorStr("Invalid ticket ID"))
		return
	}

	var body closeBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, utils.ErrorStr("Invalid request body"))
		return
	}

	// Get the ticket struct
	ticket, err := database.Client.Tickets.Get(c, ticketId, guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// Verify the ticket exists
	if ticket.UserId == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorStr("Ticket not found"))
		return
	}

	hasPermission, requestErr := utils.HasPermissionToViewTicket(context.Background(), guildId, userId, ticket)
	if requestErr != nil {
		// TODO
		c.JSON(requestErr.StatusCode, utils.ErrorJson(requestErr))
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, utils.ErrorStr("You do not have permission to close this ticket"))
		return
	}

	data := closerelay.TicketClose{
		GuildId:  guildId,
		TicketId: ticket.Id,
		UserId:   userId,
		Reason:   body.Reason,
	}

	if err := closerelay.Publish(redis.Client.Client, data); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, utils.SuccessResponse)
}

package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/closerelay"
	"github.com/gin-gonic/gin"
	"strconv"
)

type closeBody struct {
	Reason string `json:"reason"`
}

func CloseTicket(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)
	guildId := ctx.Keys["guildid"].(uint64)

	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": true,
			"error":   "Invalid ticket ID",
		})
		return
	}

	var body closeBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Missing reason",
		})
		return
	}

	// Get the ticket struct
	ticket, err := database.Client.Tickets.Get(ticketId, guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": true,
			"error":   err.Error(),
		})
		return
	}

	// Verify the ticket exists
	if ticket.UserId == 0 {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": false,
			"error":   "Ticket does not exist",
		})
		return
	}

	hasPermission, requestErr := utils.HasPermissionToViewTicket(guildId, userId, ticket)
	if err != nil {
		ctx.JSON(requestErr.StatusCode, utils.ErrorJson(requestErr))
		return
	}

	if !hasPermission {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to close this ticket"))
		return
	}

	data := closerelay.TicketClose{
		GuildId:  guildId,
		TicketId: ticket.Id,
		UserId:   userId,
		Reason:   body.Reason,
	}

	if err := closerelay.Publish(redis.Client.Client, data); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

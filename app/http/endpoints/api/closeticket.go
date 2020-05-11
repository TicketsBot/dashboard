package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/messagequeue"
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

	var data closeBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": true,
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
			"success": true,
			"error":   "Ticket does not exist",
		})
		return
	}

	if ticket.GuildId != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": true,
			"error":   "Guild ID does not matched",
		})
		return
	}

	go messagequeue.Client.PublishTicketClose(guildId, ticket.Id, userId, data.Reason)

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

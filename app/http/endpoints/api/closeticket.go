package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/gin-gonic/gin"
)

type closeBody struct {
	Reason string `json:"reason"`
}

func CloseTicket(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)
	guildId := ctx.Keys["guildid"].(uint64)
	uuid := ctx.Param("uuid")

	var data closeBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": true,
			"error":   "Missing reason",
		})
		return
	}

	// Verify that the ticket exists
	ticketChan := make(chan table.Ticket)
	go table.GetTicket(uuid, ticketChan)
	ticket := <-ticketChan

	if ticket.Uuid == "" {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": true,
			"error":   "Ticket does not exist",
		})
		return
	}

	if ticket.Guild != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": true,
			"error":   "Guild ID does not matched",
		})
		return
	}

	go messagequeue.Client.PublishTicketClose(ticket.Uuid, userId, data.Reason)

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

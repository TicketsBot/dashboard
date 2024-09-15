package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest"
	"regexp"
	"strconv"
	"time"
)

var MentionRegex, _ = regexp.Compile("<@(\\d+)>")

func GetTicket(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": true,
			"error":   "Invalid ticket ID",
		})
		return
	}

	// Get the ticket struct
	ticket, err := dbclient.Client.Tickets.Get(ctx, ticketId, guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": true,
			"error":   err.Error(),
		})
		return
	}

	if ticket.GuildId != guildId {
		ctx.JSON(403, gin.H{
			"success": false,
			"error":   "Guild ID doesn't match",
		})
		return
	}

	if !ticket.Open {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Ticket does not exist",
		})
		return
	}

	hasPermission, requestErr := utils.HasPermissionToViewTicket(context.Background(), guildId, userId, ticket)
	if requestErr != nil {
		ctx.JSON(requestErr.StatusCode, utils.ErrorJson(requestErr))
		return
	}

	if !hasPermission {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to view this ticket"))
		return
	}

	if ticket.ChannelId == nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Channel ID is nil",
		})
		return
	}

	messages, err := fetchMessages(botContext, ticket)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success":  true,
		"ticket":   ticket,
		"messages": messages,
	})
}

type StrippedMessage struct {
	Author      user.User            `json:"author"`
	Content     string               `json:"content"`
	Timestamp   time.Time            `json:"timestamp"`
	Attachments []channel.Attachment `json:"attachments"`
	Embeds      []embed.Embed        `json:"embeds"`
}

func fetchMessages(botContext *botcontext.BotContext, ticket database.Ticket) ([]StrippedMessage, error) {
	// Get messages
	messages, err := rest.GetChannelMessages(context.Background(), botContext.Token, botContext.RateLimiter, *ticket.ChannelId, rest.GetChannelMessagesData{Limit: 100})
	if err != nil {
		return nil, err
	}

	// Format messages, exclude unneeded data
	stripped := make([]StrippedMessage, len(messages))
	for i, message := range utils.Reverse(messages) {
		stripped[i] = StrippedMessage{
			Author:      message.Author,
			Content:     message.Content,
			Timestamp:   message.Timestamp,
			Attachments: message.Attachments,
			Embeds:      message.Embeds,
		}
	}

	return stripped, nil
}

package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"regexp"
	"strconv"
	"strings"
)

var MentionRegex, _ = regexp.Compile("<@(\\d+)>")

func GetTicket(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": true,
			"error":   "Invalid ticket ID",
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

	if ticket.GuildId != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error":   "Guild ID doesn't match",
		})
		return
	}

	if !ticket.Open {
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
		ctx.JSON(403, utils.ErrorStr("You do not have permission to view this ticket"))
		return
	}

	messagesFormatted := make([]map[string]interface{}, 0)
	if ticket.ChannelId != nil {
		// Get messages
		messages, _ := rest.GetChannelMessages(botContext.Token, botContext.RateLimiter, *ticket.ChannelId, rest.GetChannelMessagesData{Limit: 100})

		// Format messages, exclude unneeded data
		for _, message := range utils.Reverse(messages) {
			content := message.Content

			// Format mentions properly
			match := MentionRegex.FindAllStringSubmatch(content, -1)
			for _, mention := range match {
				if len(mention) >= 2 {
					mentionedId, err := strconv.ParseUint(mention[1], 10, 64)
					if err != nil {
						continue
					}

					user, _ := cache.Instance.GetUser(mentionedId)
					content = strings.ReplaceAll(content, fmt.Sprintf("<@%d>", mentionedId), fmt.Sprintf("@%s", user.Username))
				}
			}

			messagesFormatted = append(messagesFormatted, map[string]interface{}{
				"author":  message.Author,
				"content": content,
			})
		}
	}

	ctx.JSON(200, gin.H{
		"success":  true,
		"ticket":   ticket,
		"messages": messagesFormatted,
	})
}

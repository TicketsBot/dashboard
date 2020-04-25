package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
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
	uuid := ctx.Param("uuid")

	ticketChan := make(chan table.Ticket)
	go table.GetTicket(uuid, ticketChan)
	ticket := <-ticketChan

	if ticket.Guild != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error": "Guild ID doesn't match",
		})
		return
	}

	if !ticket.IsOpen {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": false,
			"error": "Ticket does not exist",
		})
		return
	}

	// Get messages
	messages, _ := rest.GetChannelMessages(config.Conf.Bot.Token, ratelimit.Ratelimiter, ticket.Channel, rest.GetChannelMessagesData{Limit: 100})

	// Format messages, exclude unneeded data
	messagesFormatted := make([]map[string]interface{}, 0)
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
			"username": message.Author.Username,
			"content":  content,
		})
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"ticket": ticket,
		"messages": messagesFormatted,
	})
}

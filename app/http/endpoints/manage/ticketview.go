package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"regexp"
	"strconv"
	"strings"
)

var MentionRegex, _ = regexp.Compile("<@(\\d+)>")

func TicketViewHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userId := utils.GetUserId(store)

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		guild, _ := cache.Instance.GetGuild(guildId, false)

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get ticket UUID from URL and verify it exists
		uuid := ctx.Param("uuid")
		ticketChan := make(chan table.Ticket)
		go table.GetTicket(uuid, ticketChan)
		ticket := <-ticketChan
		exists := ticket != table.Ticket{}

		// If invalid ticket UUID, take user to ticket list
		if !exists {
			ctx.Redirect(302, fmt.Sprintf("/manage/%s/tickets", guildIdStr))
			return
		}

		// Verify that the user has permission to be here
		if ticket.Guild != guildId {
			ctx.Redirect(302, fmt.Sprintf("/manage/%s/tickets", guildIdStr))
			return
		}

		// Get messages
		messages, err := rest.GetChannelMessages(config.Conf.Bot.Token, ratelimit.Ratelimiter, ticket.Channel, rest.GetChannelMessagesData{Limit: 100})

		// Format messages, exclude unneeded data
		var messagesFormatted []map[string]interface{}
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

					ch := make(chan string)
					go table.GetUsername(mentionedId, ch)
					content = strings.ReplaceAll(content, fmt.Sprintf("<@%d>", mentionedId), fmt.Sprintf("@%s", <-ch))
				}
			}

			messagesFormatted = append(messagesFormatted, map[string]interface{}{
				"username": message.Author.Username,
				"content":  content,
			})
		}

		premium := make(chan bool)
		go utils.IsPremiumGuild(store, guildId, premium)

		ctx.HTML(200, "manage/ticketview", gin.H{
			"name":         store.Get("name").(string),
			"guildId":      guildIdStr,
			"csrf":         store.Get("csrf").(string),
			"avatar":       store.Get("avatar").(string),
			"baseUrl":      config.Conf.Server.BaseUrl,
			"isError":      err != nil,
			"error":        err.Error(),
			"messages":     messagesFormatted,
			"ticketId":     ticket.TicketId,
			"uuid":         ticket.Uuid,
			"include_mock": true,
			"premium":      <-premium,
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

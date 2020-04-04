package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/channel"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SendMessage(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userIdStr := store.Get("userid").(string)
		userId, err := utils.GetUserId(store)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		var guild objects.Guild
		for _, g := range table.GetGuilds(userIdStr) {
			if g.Id == guildIdStr {
				guild = g
				break
			}
		}

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get ticket UUID from URL and verify it exists
		ticketChan := make(chan table.Ticket)
		go table.GetTicket(ctx.Param("uuid"), ticketChan)
		ticket := <-ticketChan
		exists := ticket != table.Ticket{}

		// Verify that the user has permission to be here
		if ticket.Guild != guildId {
			ctx.Redirect(302, fmt.Sprintf("/manage/%s/tickets", guildIdStr))
			return
		}

		contentType := discord.ApplicationJson

		if exists {
			content := fmt.Sprintf("**%s**: %s", store.Get("name").(string), ctx.PostForm("message"))
			if len(content) > 2000 {
				content = content[0:1999]
			}

			endpoint := channel.CreateMessage(int(ticket.Channel))
			err, _ = endpoint.Request(store, &contentType, channel.CreateMessageBody{
				Content: content,
			}, nil)
		}
	} else {
		ctx.Redirect(302, "/login")
	}

	ctx.Redirect(301, ctx.Request.URL.String())
}

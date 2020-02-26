package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/cache"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func TicketCloseHandler(ctx *gin.Context) {
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
		guildId, err := strconv.ParseInt(guildIdStr, 10, 64)
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
		if !utils.Contains(config.Conf.Admins, userIdStr) && !guild.Owner && !table.IsAdmin(guildId, userId) {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get CSRF token
		csrfCorrect := ctx.PostForm("csrf") == store.Get("csrf").(string)
		if !csrfCorrect {
			ctx.Redirect(302, "/")
			return
		}

		// Get the UUID
		uuid := ctx.Param("uuid")

		// Verify that tbe ticket exists
		ticketChan := make(chan table.Ticket)
		go table.GetTicket(uuid, ticketChan)
		ticket := <-ticketChan

		if ticket.Uuid == "" {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/tickets/view/%s?sucess=false", guildId, uuid))
			return
		}

		go cache.Client.PublishTicketClose(ticket.Uuid, userId, "") // TODO: Add option for reason

		ctx.Redirect(302, fmt.Sprintf("/manage/%d/tickets", guildId))
	} else {
		ctx.Redirect(302, "/login")
	}
}

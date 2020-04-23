package manage

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

var Archiver archiverclient.ArchiverClient

func LogViewHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}

	if utils.IsLoggedIn(store) {
		userId := utils.GetUserId(store)

		// Verify the guild exists
		guildId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		guild, _ := cache.Instance.GetGuild(guildId, false)

		// format ticket ID
		ticketId, err := strconv.Atoi(ctx.Param("ticket")); if err != nil {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/logs", guild.Id))
			return
		}

		// get ticket object
		ticketChan := make(chan table.Ticket)
		go table.GetTicketById(guildId, ticketId, ticketChan)
		ticket := <-ticketChan

		// Verify this is a valid ticket and it is closed
		if ticket.Uuid == "" || ticket.IsOpen {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/logs", guild.Id))
			return
		}

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, userId, isAdmin)
		if !<-isAdmin && ticket.Owner != userId {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// retrieve ticket messages from bucket
		messages, err := Archiver.Get(guildId, ticketId)
		if err != nil {
			if errors.Is(err, archiverclient.ErrExpired) {
				ctx.String(200, "Archives expired: Purchase premium for permanent log storage") // TODO: Actual error page
				return
			}

			ctx.String(500, fmt.Sprintf("Failed to archives - please contact the developers: %s", err.Error()))
			return
		}

		// format to html
		html, err := Archiver.Encode(messages, ticketId)
		if err != nil {
			ctx.String(500, fmt.Sprintf("Failed to archives - please contact the developers: %s", err.Error()))
			return
		}

		ctx.Data(200, gin.MIMEHTML, html)
	} else {
		ctx.Redirect(302, fmt.Sprintf("/login?noguilds&state=viewlog.%s.%s", ctx.Param("id"), ctx.Param("ticket")))
	}
}

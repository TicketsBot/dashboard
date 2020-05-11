package manage

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"strconv"
)

func ModmailLogViewHandler(ctx *gin.Context) {
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

		// get ticket UUID
		archiveUuid, err := uuid.FromString(ctx.Param("uuid"))
		if err != nil {
			// TODO: 404 error page
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"error": "Modmail archive not found",
			})
			return
		}

		// get ticket object
		archive, err := database.Client.ModmailArchive.Get(archiveUuid)
		if err != nil {
			// TODO: 500 error page
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
			return
		}

		// Verify this is a valid ticket and it is closed
		if archive.Uuid == uuid.Nil{
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/logs/modmail", guild.Id))
			return
		}

		// Verify this modmail ticket was for this guild
		if archive.GuildId != guildId {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/logs/modmail", guild.Id))
			return
		}

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, userId, isAdmin)
		if !<-isAdmin && archive.UserId != userId {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// retrieve ticket messages from bucket
		messages, err := Archiver.GetModmail(guildId, archiveUuid.String())
		if err != nil {
			if errors.Is(err, archiverclient.ErrExpired) {
				ctx.String(200, "Archives expired: Purchase premium for permanent log storage") // TODO: Actual error page
				return
			}

			ctx.String(500, fmt.Sprintf("Failed to retrieve archive - please contact the developers: %s", err.Error()))
			return
		}

		// format to html
		html, err := Archiver.Encode(messages, fmt.Sprintf("modmail-%s", archiveUuid))
		if err != nil {
			ctx.String(500, fmt.Sprintf("Failed to retrieve archive - please contact the developers: %s", err.Error()))
			return
		}

		ctx.Data(200, gin.MIMEHTML, html)
	} else {
		ctx.Redirect(302, fmt.Sprintf("/login?noguilds&state=viewlog.%s.%s", ctx.Param("id"), ctx.Param("ticket")))
	}
}

package manage

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/TicketsBot/common/permission"
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
			utils.ErrorPage(ctx, 404, "Couldn't find a server with the ID provided")
			return
		}

		// get ticket UUID
		archiveUuid, err := uuid.FromString(ctx.Param("uuid"))
		if err != nil {
			utils.ErrorPage(ctx, 404, "Modmail archive with provided UUID not found")
			return
		}

		// get ticket object
		archive, err := database.Client.ModmailArchive.Get(archiveUuid)
		if err != nil {
			utils.ErrorPage(ctx, 500, err.Error())
			return
		}

		// Verify this is a valid ticket and it is closed
		if archive.Uuid == uuid.Nil {
			utils.ErrorPage(ctx, 404, "Modmail archive with provided UUID not found")
			return
		}

		// Verify this modmail ticket was for this guild
		if archive.GuildId != guildId {
			utils.ErrorPage(ctx, 403, "Modmail archive does not belong to this server")
			return
		}

		// Verify the user has permissions to be here
		permLevel, err := utils.GetPermissionLevel(guildId, userId)
		if err != nil {
			ctx.JSON(500, utils.ErrorToResponse(err))
			return
		}

		if permLevel < permission.Support && archive.UserId != userId {
			utils.ErrorPage(ctx, 403, "You do not have permission to view this archive")
			return
		}

		// retrieve ticket messages from bucket
		messages, err := Archiver.GetModmail(guildId, archiveUuid.String())
		if err != nil {
			if errors.Is(err, archiverclient.ErrExpired) {
				utils.ErrorPage(ctx, 404, "Archive expired - purchase premium for permanent log storage")
			} else {
				utils.ErrorPage(ctx, 500, err.Error())
			}

			return
		}

		// format to html
		html, err := Archiver.Encode(messages, fmt.Sprintf("modmail-%s", archiveUuid))
		if err != nil {
			utils.ErrorPage(ctx, 500, err.Error())
			return
		}

		ctx.Data(200, gin.MIMEHTML, html)
	} else {
		ctx.Redirect(302, fmt.Sprintf("/login?noguilds&state=viewlog.%s.%s", ctx.Param("id"), ctx.Param("ticket")))
	}
}

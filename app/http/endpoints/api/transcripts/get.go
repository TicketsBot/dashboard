package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetTranscriptHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	// format ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid ticket ID"))
		return
	}

	// get ticket object
	ticket, err := database.Client.Tickets.Get(ticketId, guildId)
	if err != nil {
		// TODO: 500 error page
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Verify this is a valid ticket and it is closed
	if ticket.UserId == 0 || ticket.Open {
		ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
	if ticket.UserId != userId {
		permLevel, err := utils.GetPermissionLevel(guildId, userId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if permLevel < permission.Support {
			ctx.JSON(403, utils.ErrorStr("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	messages, err := utils.ArchiverClient.Get(guildId, ticketId)
	if err != nil {
		if errors.Is(err, archiverclient.ErrExpired) {
			ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	ctx.JSON(200, messages)
}

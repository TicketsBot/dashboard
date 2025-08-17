package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/chatreplica"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetTranscriptRenderHandler(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)
	userId := c.Keys["userid"].(uint64)

	// format ticket ID
	ticketId, err := strconv.Atoi(c.Param("ticketId"))
	if err != nil {
		c.JSON(400, utils.ErrorStr("Invalid ticket ID"))
		return
	}

	// get ticket object
	ticket, err := dbclient.Client.Tickets.Get(c, ticketId, guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// Verify this is a valid ticket and it is closed
	if ticket.UserId == 0 || ticket.Open {
		c.JSON(404, utils.ErrorStr("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
	// ticket.UserId cannot be 0
	if ticket.UserId != userId {
		hasPermission, err := utils.HasPermissionToViewTicket(c, guildId, userId, ticket)
		if err != nil {
			if err.StatusCode >= 500 {
				_ = c.AbortWithError(err.StatusCode, app.NewServerError(err))
			} else {
				c.JSON(err.StatusCode, utils.ErrorJson(err))
			}
			return
		}

		if !hasPermission {
			c.JSON(403, utils.ErrorStr("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	transcript, err := utils.ArchiverClient.Get(c, guildId, ticketId)
	if err != nil {
		if errors.Is(err, archiverclient.ErrNotFound) {
			c.JSON(404, utils.ErrorStr("Transcript not found"))
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		}

		return
	}

	// Render
	payload := chatreplica.FromTranscript(transcript, ticketId)
	html, err := chatreplica.Render(payload)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.Data(200, "text/html", html)
}

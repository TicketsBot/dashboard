package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"net/http"
	"strconv"
)

func MultiPanelDelete(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	multiPanelId, err := strconv.Atoi(c.Param("panelid"))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	panel, ok, err := dbclient.Client.MultiPanels.Get(c, multiPanelId)
	if !ok {
		c.JSON(404, utils.ErrorStr("No panel with matching ID found"))
		return
	}

	if panel.GuildId != guildId {
		c.JSON(403, utils.ErrorStr("Guild ID doesn't match"))
		return
	}

	// TODO: Use proper context
	if err := rest.DeleteMessage(c, botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.MessageId); err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) {
			// Swallow 403 / 404
			if unwrapped.StatusCode != http.StatusForbidden && unwrapped.StatusCode != http.StatusNotFound {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			return
		}
	}

	success, err := dbclient.Client.MultiPanels.Delete(c, guildId, multiPanelId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	if !success {
		c.JSON(404, utils.ErrorJson(errors.New("No panel with matching ID found")))
		return
	}

	c.JSON(200, utils.SuccessResponse)
}

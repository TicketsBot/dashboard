package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

func MultiPanelDelete(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	multiPanelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panel, ok, err := dbclient.Client.MultiPanels.Get(multiPanelId)
	if !ok {
		ctx.JSON(404, utils.ErrorStr("No panel with matching ID found"))
		return
	}

	if panel.GuildId != guildId {
		ctx.JSON(403, utils.ErrorStr("Guild ID doesn't match"))
		return
	}

	var unwrapped request.RestError
	// TODO: Use proper context
	if err := rest.DeleteMessage(context.Background(), botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.MessageId); err != nil && !(errors.As(err, &unwrapped) && unwrapped.IsClientError()) {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	success, err := dbclient.Client.MultiPanels.Delete(guildId, multiPanelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !success {
		ctx.JSON(404, utils.ErrorJson(errors.New("No panel with matching ID found")))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

func MultiPanelResend(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// parse panel ID
	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// retrieve panel from DB
	multiPanel, ok, err := dbclient.Client.MultiPanels.Get(panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// check panel exists
	if !ok {
		ctx.JSON(404, utils.ErrorStr("No panel with the provided ID found"))
		return
	}

	// check panel is in the same guild
	if guildId != multiPanel.GuildId {
		ctx.JSON(403, utils.ErrorStr("Guild ID doesn't match"))
		return
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// delete old message
	if err := rest.DeleteMessage(botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId); err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && !unwrapped.IsClientError() {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	// get premium status
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panels, err := dbclient.Client.MultiPanelTargets.GetPanels(multiPanel.Id)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// send new message
	messageData := multiPanelIntoMessageData(multiPanel, premiumTier > premium.None)
	messageId, err := messageData.send(&botContext, panels)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			ctx.JSON(500, utils.ErrorJson(errors.New("I do not have permission to send messages in the provided channel")))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	if err = dbclient.Client.MultiPanels.UpdateMessageId(multiPanel.Id, messageId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

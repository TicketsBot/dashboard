package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func MultiPanelUpdate(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// parse body
	var data multiPanelCreateData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

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
		ctx.JSON(404, utils.ErrorJson(errors.New("No panel with the provided ID found")))
		return
	}

	// check panel is in the same guild
	if guildId != multiPanel.GuildId {
		ctx.JSON(403, utils.ErrorJson(errors.New("Guild ID doesn't match")))
		return
	}

	// validate body & get sub-panels
	panels, err := data.doValidations(guildId)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	for _, panel := range panels {
		if panel.CustomId == "" {
			panel.CustomId = utils.RandString(80)
			if err := dbclient.Client.Panel.Update(panel); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		}
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// delete old message
	var unwrapped request.RestError
	if err := rest.DeleteMessage(botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId); err != nil && !(errors.As(err, &unwrapped) && unwrapped.IsClientError())  {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// get premium status
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// send new message
	messageData := data.IntoMessageData(premiumTier > premium.None)
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

	// update DB
	updated := database.MultiPanel{
		Id:        multiPanel.Id,
		MessageId: messageId,
		ChannelId: data.ChannelId,
		GuildId:   guildId,
		Title:     data.Title,
		Content:   data.Content,
		Colour:    int(data.Colour),
	}

	if err = dbclient.Client.MultiPanels.Update(multiPanel.Id, updated); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// TODO: one query for ACID purposes
	// delete old targets
	if err := dbclient.Client.MultiPanelTargets.DeleteAll(multiPanel.Id); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// insert new targets
	group, _ := errgroup.WithContext(context.Background())
	for _, panel := range panels {
		panel := panel

		group.Go(func() error {
			return dbclient.Client.MultiPanelTargets.Insert(multiPanel.Id, panel.PanelId)
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    multiPanel,
	})
}

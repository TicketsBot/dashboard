package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

func DeletePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	panel, err := database.Client.Panel.GetById(panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if panel.PanelId == 0 {
		ctx.JSON(404, utils.ErrorStr("Panel not found"))
		return
	}

	// verify panel belongs to guild
	if panel.GuildId != guildId {
		ctx.JSON(403, utils.ErrorStr("Guild ID doesn't match"))
		return
	}

	// Get any multi panels this panel is part of to use later
	multiPanels, err := database.Client.MultiPanelTargets.GetMultiPanels(panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Delete welcome message embed
	if panel.WelcomeMessageEmbed != nil {
		if err := database.Client.Embeds.Delete(*panel.WelcomeMessageEmbed); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	if err := database.Client.Panel.Delete(panelId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if err := rest.DeleteMessage(botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.MessageId); err != nil {
		var unwrapped request.RestError
		if !errors.As(err, &unwrapped) || unwrapped.StatusCode != 404 {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	// Get premium tier
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Update all multi panels messages to remove the button
	for i, multiPanel := range multiPanels {
		// Only update 5 multi-panels maximum: Prevent DoS
		if i >= 5 {
			break
		}

		panels, err := database.Client.MultiPanelTargets.GetPanels(multiPanel.Id)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		messageData := multiPanelMessageData{
			Title:      multiPanel.Title,
			Content:    multiPanel.Content,
			Colour:     multiPanel.Colour,
			ChannelId:  multiPanel.ChannelId,
			SelectMenu: multiPanel.SelectMenu,
			IsPremium:  premiumTier > premium.None,
		}

		messageId, err := messageData.send(&botContext, panels)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if err := database.Client.MultiPanels.UpdateMessageId(multiPanel.Id, messageId); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Delete old panel
		_ = rest.DeleteMessage(botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId)
	}

	ctx.JSON(200, utils.SuccessResponse)
}

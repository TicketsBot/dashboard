package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

func UpdatePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	var data panelBody

	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, utils.ErrorJson(err))
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, utils.ErrorJson(err))
		return
	}

	// get existing
	existing, err := dbclient.Client.Panel.GetById(panelId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// check guild ID matches
	if existing.GuildId != guildId {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Guild ID does not match",
		})
		return
	}

	if !data.doValidations(ctx, guildId) {
		return
	}

	// check if this will break a multi-panel;
	// first, get any multipanels this panel belongs to
	multiPanels, err := dbclient.Client.MultiPanelTargets.GetMultiPanels(existing.PanelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	premiumTier := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)

	for _, multiPanel := range multiPanels {
		panels, err := dbclient.Client.MultiPanelTargets.GetPanels(multiPanel.Id)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// TODO: Optimise this
		panelIds := make([]int, len(panels))
		for i, panel := range panels {
			panelIds[i] = panel.PanelId
		}

		data := multiPanelCreateData{
			Title:     multiPanel.Title,
			Content:   multiPanel.Content,
			Colour:    int32(multiPanel.Colour),
			ChannelId: multiPanel.ChannelId,
			Panels:    panelIds,
		}

		messageId, err := data.sendEmbed(&botContext, premiumTier > premium.None, panels)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if err := dbclient.Client.MultiPanels.UpdateMessageId(multiPanel.Id, messageId); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	// check if we need to update the message
	shouldUpdateMessage := uint32(existing.Colour) != data.Colour ||
		existing.ChannelId != data.ChannelId ||
		existing.Content != data.Content ||
		existing.Title != data.Title ||
		existing.ReactionEmote != data.Emote ||
		existing.ImageUrl != data.ImageUrl ||
		existing.ThumbnailUrl != data.ThumbnailUrl ||
		component.ButtonStyle(existing.ButtonStyle) != data.ButtonStyle

	emoji, _ := data.getEmoji() // already validated
	newMessageId := existing.MessageId

	if shouldUpdateMessage {
		// delete old message, ignoring error
		_ = rest.DeleteMessage(botContext.Token, botContext.RateLimiter, existing.ChannelId, existing.MessageId)

		newMessageId, err = data.sendEmbed(&botContext, data.Title, existing.CustomId, data.Emote, data.ImageUrl, data.ThumbnailUrl, data.ButtonStyle, premiumTier > premium.None)
		if err != nil {
			var unwrapped request.RestError
			if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
				ctx.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   "I do not have permission to send messages in the specified channel",
				})
			} else {
				// TODO: Most appropriate error?
				ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
			}

			return
		}
	}

	// Store in DB
	panel := database.Panel{
		PanelId:         panelId,
		MessageId:       newMessageId,
		ChannelId:       data.ChannelId,
		GuildId:         guildId,
		Title:           data.Title,
		Content:         data.Content,
		Colour:          int32(data.Colour),
		TargetCategory:  data.CategoryId,
		ReactionEmote:   emoji,
		WelcomeMessage:  data.WelcomeMessage,
		WithDefaultTeam: data.WithDefaultTeam,
		CustomId:        existing.CustomId,
		ImageUrl:        data.ImageUrl,
		ThumbnailUrl:    data.ThumbnailUrl,
		ButtonStyle:     int(data.ButtonStyle),
	}

	if err = dbclient.Client.Panel.Update(panel); err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// insert role mention data
	// delete old data
	if err = dbclient.Client.PanelRoleMentions.DeleteAll(panel.PanelId); err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// string is role ID or "user" to mention the ticket opener
	var shouldMentionUser bool
	for _, mention := range data.Mentions {
		if mention == "user" {
			shouldMentionUser = true
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
				return
			}

			// should we check the role is a valid role in the guild?
			// not too much of an issue if it isnt
			if err = dbclient.Client.PanelRoleMentions.Add(panel.PanelId, roleId); err != nil {
				ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
				return
			}
		}
	}

	if err = dbclient.Client.PanelUserMention.Set(panel.PanelId, shouldMentionUser); err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// insert support teams
	// TODO: Stop race conditions - 1 transaction
	// delete teams
	if err := dbclient.Client.PanelTeams.DeleteAll(panel.PanelId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// insert new
	if responseCode, err := insertTeams(guildId, panel.PanelId, data.Teams); err != nil {
		ctx.JSON(responseCode, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

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
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	var data panelBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get existing
	existing, err := dbclient.Client.Panel.GetById(panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// check guild ID matches
	if existing.GuildId != guildId {
		ctx.JSON(400, utils.ErrorStr("Guild ID does not match"))
		return
	}

	if !data.doValidations(ctx, guildId) {
		return
	}

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	var emojiId *uint64
	var emojiName *string
	{
		emoji := data.getEmoji()
		if emoji != nil {
			emojiName = &emoji.Name

			if emoji.Id.Value != 0 {
				emojiId = &emoji.Id.Value
			}
		}
	}

	// check if we need to update the message
	shouldUpdateMessage := uint32(existing.Colour) != data.Colour ||
		existing.ChannelId != data.ChannelId ||
		existing.Content != data.Content ||
		existing.Title != data.Title ||
		(existing.EmojiId == nil && emojiId != nil || existing.EmojiId != nil && emojiId == nil || (existing.EmojiId != nil && emojiId != nil && *existing.EmojiId != *emojiId)) ||
		(existing.EmojiName == nil && emojiName != nil || existing.EmojiName != nil && emojiName == nil || (existing.EmojiName != nil && emojiName != nil && *existing.EmojiName != *emojiName)) ||
		existing.ImageUrl != data.ImageUrl ||
		existing.ThumbnailUrl != data.ThumbnailUrl ||
		component.ButtonStyle(existing.ButtonStyle) != data.ButtonStyle ||
		existing.ButtonLabel != data.ButtonLabel

	newMessageId := existing.MessageId

	if shouldUpdateMessage {
		// delete old message, ignoring error
		_ = rest.DeleteMessage(botContext.Token, botContext.RateLimiter, existing.ChannelId, existing.MessageId)

		messageData := data.IntoPanelMessageData(existing.CustomId, premiumTier > premium.None)
		newMessageId, err = messageData.send(&botContext)
		if err != nil {
			var unwrapped request.RestError
			if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
				ctx.JSON(500, utils.ErrorStr("I do not have permission to send messages in the specified channel"))
			} else {
				// TODO: Most appropriate error?
				ctx.JSON(500, utils.ErrorJson(err))
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
		EmojiName:       emojiName,
		EmojiId:         emojiId,
		WelcomeMessage:  data.WelcomeMessage,
		WithDefaultTeam: data.WithDefaultTeam,
		CustomId:        existing.CustomId,
		ImageUrl:        data.ImageUrl,
		ThumbnailUrl:    data.ThumbnailUrl,
		ButtonStyle:     int(data.ButtonStyle),
		ButtonLabel:     data.ButtonLabel,
		FormId:          data.FormId,
	}

	if err = dbclient.Client.Panel.Update(panel); err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	// insert mention data
	validRoles, err := getRoleHashSet(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// string is role ID or "user" to mention the ticket opener
	var shouldMentionUser bool
	var roleMentions []uint64
	for _, mention := range data.Mentions {
		if mention == "user" {
			shouldMentionUser = true
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
				return
			}

			if validRoles.Contains(roleId) {
				roleMentions = append(roleMentions, roleId)
			}
		}
	}

	if err := dbclient.Client.PanelUserMention.Set(panel.PanelId, shouldMentionUser); err != nil {
		ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
		return
	}

	if err := dbclient.Client.PanelRoleMentions.Replace(panel.PanelId, roleMentions); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// We are safe to insert, team IDs already validated
	if err := dbclient.Client.PanelTeams.Replace(panel.PanelId, data.Teams); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Update multi panels

	// check if this will break a multi-panel;
	// first, get any multipanels this panel belongs to
	multiPanels, err := dbclient.Client.MultiPanelTargets.GetMultiPanels(existing.PanelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	for i, multiPanel := range multiPanels {
		// Only update 5 multi-panels maximum: Prevent DoS
		if i >= 5 {
			break
		}

		panels, err := dbclient.Client.MultiPanelTargets.GetPanels(multiPanel.Id)
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

		if err := dbclient.Client.MultiPanels.UpdateMessageId(multiPanel.Id, messageId); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Delete old panel
		_ = rest.DeleteMessage(botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId)
	}

	ctx.JSON(200, utils.SuccessResponse)
}

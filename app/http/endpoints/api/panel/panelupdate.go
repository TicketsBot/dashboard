package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/app/http/validation"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	if existing.ForceDisabled {
		ctx.JSON(400, utils.ErrorStr("This panel is disabled and cannot be modified: please reactivate premium to re-enable it"))
		return
	}

	// Apply defaults
	ApplyPanelDefaults(&data)

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	channels, err := botContext.GetGuildChannels(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Do custom validation
	validationContext := PanelValidationContext{
		Data:       data,
		GuildId:    guildId,
		IsPremium:  premiumTier > premium.None,
		BotContext: botContext,
		Channels:   channels,
	}

	if err := ValidatePanelBody(validationContext); err != nil {
		var validationError *validation.InvalidInputError
		if errors.As(err, &validationError) {
			ctx.JSON(400, utils.ErrorStr(validationError.Error()))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	// Do tag validation
	if err := validate.Struct(data); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(500, utils.ErrorStr("An error occurred while validating the panel"))
			return
		}

		formatted := "Your input contained the following errors:\n" + utils.FormatValidationErrors(validationErrors)
		ctx.JSON(400, utils.ErrorStr(formatted))
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
		existing.ButtonLabel != data.ButtonLabel ||
		existing.Disabled != data.Disabled

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

	// Update welcome message
	var welcomeMessageEmbed *int
	if data.WelcomeMessage == nil {
		if existing.WelcomeMessageEmbed != nil { // If welcome message wasn't null, but now is, delete the embed
			if err := dbclient.Client.Embeds.Delete(*existing.WelcomeMessageEmbed); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		} // else, welcomeMessageEmbed will be nil
	} else {
		// TODO: Upsert? Don't think we can, as no unique key in the table, panel_id is in panels table
		if existing.WelcomeMessageEmbed == nil { // Create
			embed, fields := data.WelcomeMessage.IntoDatabaseStruct()
			embed.GuildId = guildId

			id, err := dbclient.Client.Embeds.CreateWithFields(embed, fields)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			welcomeMessageEmbed = &id
		} else { // Update
			welcomeMessageEmbed = existing.WelcomeMessageEmbed

			embed, fields := data.WelcomeMessage.IntoDatabaseStruct()
			embed.Id = *existing.WelcomeMessageEmbed
			embed.GuildId = guildId

			if err := dbclient.Client.Embeds.UpdateWithFields(embed, fields); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		}
	}

	// Store in DB
	panel := database.Panel{
		PanelId:             panelId,
		MessageId:           newMessageId,
		ChannelId:           data.ChannelId,
		GuildId:             guildId,
		Title:               data.Title,
		Content:             data.Content,
		Colour:              int32(data.Colour),
		TargetCategory:      data.CategoryId,
		EmojiName:           emojiName,
		EmojiId:             emojiId,
		WelcomeMessageEmbed: welcomeMessageEmbed,
		WithDefaultTeam:     data.WithDefaultTeam,
		CustomId:            existing.CustomId,
		ImageUrl:            data.ImageUrl,
		ThumbnailUrl:        data.ThumbnailUrl,
		ButtonStyle:         int(data.ButtonStyle),
		ButtonLabel:         data.ButtonLabel,
		FormId:              data.FormId,
		NamingScheme:        data.NamingScheme,
		ForceDisabled:       existing.ForceDisabled,
		Disabled:            data.Disabled,
		ExitSurveyFormId:    data.ExitSurveyFormId,
	}

	if err = dbclient.Client.Panel.Update(panel); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
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
				ctx.JSON(500, utils.ErrorJson(err))
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
			Title:        multiPanel.Title,
			Content:      multiPanel.Content,
			Colour:       multiPanel.Colour,
			ChannelId:    multiPanel.ChannelId,
			SelectMenu:   multiPanel.SelectMenu,
			IsPremium:    premiumTier > premium.None,
			ImageUrl:     multiPanel.ImageUrl,
			ThumbnailUrl: multiPanel.ThumbnailUrl,
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

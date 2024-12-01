package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/app/http/validation"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"net/http"
	"strconv"
)

func UpdatePanel(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	var data panelBody
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorStr("Invalid request body"))
		return
	}

	panelId, err := strconv.Atoi(c.Param("panelid"))
	if err != nil {
		c.JSON(400, utils.ErrorStr("Missing panel ID"))
		return
	}

	// get existing
	existing, err := dbclient.Client.Panel.GetById(c, panelId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// check guild ID matches
	if existing.GuildId != guildId {
		c.JSON(400, utils.ErrorStr("Guild ID does not match"))
		return
	}

	if existing.ForceDisabled {
		c.JSON(400, utils.ErrorStr("This panel is disabled and cannot be modified: please reactivate premium to re-enable it"))
		return
	}

	// Apply defaults
	ApplyPanelDefaults(&data)

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(c, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// TODO: Use proper context
	channels, err := botContext.GetGuildChannels(context.Background(), guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// TODO: Use proper context
	roles, err := botContext.GetGuildRoles(context.Background(), guildId)
	if err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// Do custom validation
	validationContext := PanelValidationContext{
		Data:       data,
		GuildId:    guildId,
		IsPremium:  premiumTier > premium.None,
		BotContext: botContext,
		Channels:   channels,
		Roles:      roles,
	}

	if err := ValidatePanelBody(validationContext); err != nil {
		var validationError *validation.InvalidInputError
		if errors.As(err, &validationError) {
			c.JSON(400, utils.ErrorStr(validationError.Error()))
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		}

		return
	}

	// Do tag validation
	if err := validate.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			c.JSON(500, utils.ErrorStr("An error occurred while validating the panel"))
			return
		}

		formatted := "Your input contained the following errors:\n" + utils.FormatValidationErrors(validationErrors)
		c.JSON(400, utils.ErrorStr(formatted))
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
		// TODO: Use proper context
		_ = rest.DeleteMessage(c, botContext.Token, botContext.RateLimiter, existing.ChannelId, existing.MessageId)

		messageData := data.IntoPanelMessageData(existing.CustomId, premiumTier > premium.None)
		newMessageId, err = messageData.send(botContext)
		if err != nil {
			var unwrapped request.RestError
			if errors.As(err, &unwrapped) {
				if unwrapped.StatusCode == 403 {
					c.JSON(403, utils.ErrorStr("I do not have permission to send messages in the specified channel"))
					return
				} else if unwrapped.StatusCode == 404 {
					// Swallow error
					// TODO: Make channel_id column nullable, and set to null
				} else {
					_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
					return
				}
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}
		}
	}

	// Update welcome message
	var welcomeMessageEmbed *int
	if data.WelcomeMessage == nil {
		if existing.WelcomeMessageEmbed != nil { // If welcome message wasn't null, but now is, delete the embed
			if err := dbclient.Client.Embeds.Delete(c, *existing.WelcomeMessageEmbed); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}
		} // else, welcomeMessageEmbed will be nil
	} else {
		// TODO: Upsert? Don't think we can, as no unique key in the table, panel_id is in panels table
		if existing.WelcomeMessageEmbed == nil { // Create
			embed, fields := data.WelcomeMessage.IntoDatabaseStruct()
			embed.GuildId = guildId

			id, err := dbclient.Client.Embeds.CreateWithFields(c, embed, fields)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}

			welcomeMessageEmbed = &id
		} else { // Update
			welcomeMessageEmbed = existing.WelcomeMessageEmbed

			embed, fields := data.WelcomeMessage.IntoDatabaseStruct()
			embed.Id = *existing.WelcomeMessageEmbed
			embed.GuildId = guildId

			if err := dbclient.Client.Embeds.UpdateWithFields(c, embed, fields); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
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
		PendingCategory:     data.PendingCategory,
	}

	// insert mention data
	validRoles := utils.ToSet(utils.Map(roles, utils.RoleToId))

	// string is role ID or "user" to mention the ticket opener
	var shouldMentionUser bool
	var roleMentions []uint64
	for _, mention := range data.Mentions {
		if mention == "user" {
			shouldMentionUser = true
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}

			if validRoles.Contains(roleId) {
				roleMentions = append(roleMentions, roleId)
			}
		}
	}

	err = dbclient.Client.Panel.BeginFunc(c, func(tx pgx.Tx) error {
		if err := dbclient.Client.Panel.UpdateWithTx(c, tx, panel); err != nil {
			return err
		}

		if err := dbclient.Client.PanelUserMention.SetWithTx(c, tx, panel.PanelId, shouldMentionUser); err != nil {
			return err
		}

		if err := dbclient.Client.PanelRoleMentions.ReplaceWithTx(c, tx, panel.PanelId, roleMentions); err != nil {
			return err
		}

		// We are safe to insert, team IDs already validated
		if err := dbclient.Client.PanelTeams.ReplaceWithTx(c, tx, panel.PanelId, data.Teams); err != nil {
			return err
		}

		if err := dbclient.Client.PanelAccessControlRules.ReplaceWithTx(c, tx, panel.PanelId, data.AccessControlList); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// This doesn't need to be done in a transaction
	// Update multi panels

	// check if this will break a multi-panel;
	// first, get any multipanels this panel belongs to
	multiPanels, err := dbclient.Client.MultiPanelTargets.GetMultiPanels(c, existing.PanelId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	for i, multiPanel := range multiPanels {
		// Only update 5 multi-panels maximum: Prevent DoS
		if i >= 5 {
			break
		}

		panels, err := dbclient.Client.MultiPanelTargets.GetPanels(c, multiPanel.Id)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			return
		}

		messageData := multiPanelIntoMessageData(multiPanel, premiumTier > premium.None)

		messageId, err := messageData.send(botContext, panels)
		if err != nil {
			var unwrapped request.RestError
			if errors.As(err, &unwrapped) {
				if unwrapped.StatusCode == http.StatusForbidden {
					c.JSON(400, utils.ErrorStr("I do not have permission to send messages in the specified channel"))
				} else {
					c.JSON(400, utils.ErrorStr("Error sending panel message: "+unwrapped.ApiError.Message))
				}
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			}

			return
		}

		if err := dbclient.Client.MultiPanels.UpdateMessageId(c, multiPanel.Id, messageId); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			return
		}

		// Delete old panel
		// TODO: Use proper context
		_ = rest.DeleteMessage(c, botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId)
	}

	c.JSON(200, utils.SuccessResponse)
}

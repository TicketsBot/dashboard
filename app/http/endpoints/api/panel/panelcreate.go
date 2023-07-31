package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/app/http/validation"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/common/collections"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

const freePanelLimit = 3

type panelBody struct {
	ChannelId        uint64                `json:"channel_id,string"`
	MessageId        uint64                `json:"message_id,string"`
	Title            string                `json:"title"`
	Content          string                `json:"content"`
	Colour           uint32                `json:"colour"`
	CategoryId       uint64                `json:"category_id,string"`
	Emoji            types.Emoji           `json:"emote"`
	WelcomeMessage   *types.CustomEmbed    `json:"welcome_message" validate:"omitempty,dive"`
	Mentions         []string              `json:"mentions"`
	WithDefaultTeam  bool                  `json:"default_team"`
	Teams            []int                 `json:"teams"`
	ImageUrl         *string               `json:"image_url,omitempty"`
	ThumbnailUrl     *string               `json:"thumbnail_url,omitempty"`
	ButtonStyle      component.ButtonStyle `json:"button_style,string"`
	ButtonLabel      string                `json:"button_label"`
	FormId           *int                  `json:"form_id"`
	NamingScheme     *string               `json:"naming_scheme"`
	Disabled         bool                  `json:"disabled"`
	ExitSurveyFormId *int                  `json:"exit_survey_form_id"`
}

func (p *panelBody) IntoPanelMessageData(customId string, isPremium bool) panelMessageData {
	return panelMessageData{
		ChannelId:      p.ChannelId,
		Title:          p.Title,
		Content:        p.Content,
		CustomId:       customId,
		Colour:         int(p.Colour),
		ImageUrl:       p.ImageUrl,
		ThumbnailUrl:   p.ThumbnailUrl,
		Emoji:          p.getEmoji(),
		ButtonStyle:    p.ButtonStyle,
		ButtonLabel:    p.ButtonLabel,
		ButtonDisabled: p.Disabled,
		IsPremium:      isPremium,
	}
}

var validate = validator.New()

func CreatePanel(ctx *gin.Context) {
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

	data.MessageId = 0

	// Check panel quota
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, false, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if premiumTier == premium.None {
		panels, err := dbclient.Client.Panel.GetByGuild(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if len(panels) >= freePanelLimit {
			ctx.JSON(402, utils.ErrorStr("You have exceeded your panel quota. Purchase premium to unlock more panels."))
			return
		}
	}

	// Apply defaults
	ApplyPanelDefaults(&data)

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

	customId, err := utils.RandString(30)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	messageData := data.IntoPanelMessageData(customId, premiumTier > premium.None)
	msgId, err := messageData.send(&botContext)
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

	// Store welcome message embed first
	var welcomeMessageEmbed *int
	if data.WelcomeMessage != nil {
		embed, fields := data.WelcomeMessage.IntoDatabaseStruct()
		embed.GuildId = guildId

		id, err := dbclient.Client.Embeds.CreateWithFields(embed, fields)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		welcomeMessageEmbed = &id
	}

	// Store in DB
	panel := database.Panel{
		MessageId:           msgId,
		ChannelId:           data.ChannelId,
		GuildId:             guildId,
		Title:               data.Title,
		Content:             data.Content,
		Colour:              int32(data.Colour),
		TargetCategory:      data.CategoryId,
		EmojiId:             emojiId,
		EmojiName:           emojiName,
		WelcomeMessageEmbed: welcomeMessageEmbed,
		WithDefaultTeam:     data.WithDefaultTeam,
		CustomId:            customId,
		ImageUrl:            data.ImageUrl,
		ThumbnailUrl:        data.ThumbnailUrl,
		ButtonStyle:         int(data.ButtonStyle),
		ButtonLabel:         data.ButtonLabel,
		FormId:              data.FormId,
		NamingScheme:        data.NamingScheme,
		ForceDisabled:       false,
		Disabled:            data.Disabled,
		ExitSurveyFormId:    data.ExitSurveyFormId,
	}

	panelId, err := dbclient.Client.Panel.Create(panel)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// insert role mention data
	// string is role ID or "user" to mention the ticket opener
	validRoles, err := getRoleHashSet(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	var roleMentions []uint64
	for _, mention := range data.Mentions {
		if mention == "user" {
			if err = dbclient.Client.PanelUserMention.Set(panelId, true); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				ctx.JSON(400, utils.ErrorStr("Invalid role ID"))
				return
			}

			if validRoles.Contains(roleId) {
				roleMentions = append(roleMentions, roleId)
			}
		}
	}

	if err := dbclient.Client.PanelRoleMentions.Replace(panelId, roleMentions); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Already validated, we are safe to insert
	if err := dbclient.Client.PanelTeams.Replace(panelId, data.Teams); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success":  true,
		"panel_id": panelId,
	})
}

// Data must be validated before calling this function
func (p *panelBody) getEmoji() *emoji.Emoji {
	return p.Emoji.IntoGdl()
}

func getRoleHashSet(guildId uint64) (*collections.Set[uint64], error) {
	ctx, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return nil, err
	}

	roles, err := ctx.GetGuildRoles(guildId)
	if err != nil {
		return nil, err
	}

	set := collections.NewSet[uint64]()

	for _, role := range roles {
		set.Add(role.Id)
	}

	return set, nil
}

package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/app/http/validation"
	"github.com/TicketsBot/GoPanel/app/http/validation/defaults"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/interaction/component"
	"regexp"
	"strings"
	"time"
)

func ApplyPanelDefaults(data *panelBody) {
	for _, applicator := range DefaultApplicators(data) {
		if applicator.ShouldApply() {
			applicator.Apply()
		}
	}
}

func DefaultApplicators(data *panelBody) []defaults.DefaultApplicator {
	return []defaults.DefaultApplicator{
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.Title, "Open a ticket!"),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.Content, "By clicking the button, a ticket will be opened for you."),
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.ImageUrl, nil),
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.ThumbnailUrl, nil),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.ButtonLabel, data.Title),
		defaults.NewDefaultApplicator(defaults.EmptyStringCheck, &data.ButtonLabel, "Open a ticket!"), // Title could have been blank
		defaults.NewDefaultApplicator[*string](defaults.NilOrEmptyStringCheck, &data.NamingScheme, nil),
	}
}

type PanelValidationContext struct {
	Data       panelBody
	GuildId    uint64
	IsPremium  bool
	BotContext *botcontext.BotContext
	Channels   []channel.Channel
	Roles      []guild.Role
}

func ValidatePanelBody(validationContext PanelValidationContext) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	return validation.Validate(ctx, validationContext, panelValidators()...)
}

func panelValidators() []validation.Validator[PanelValidationContext] {
	return []validation.Validator[PanelValidationContext]{
		validateTitle,
		validateContent,
		validateChannelId,
		validateCategory,
		validateEmoji,
		validateImageUrl,
		validateThumbnailUrl,
		validateButtonStyle,
		validateButtonLabel,
		validateFormId,
		validateExitSurveyFormId,
		validateTeams,
		validateNamingScheme,
		validateWelcomeMessage,
		validateAccessControlList,
	}
}

func validateTitle(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.Title) > 80 {
			return validation.NewInvalidInputError("Panel title must be less than 80 characters")
		}

		return nil
	}
}

func validateContent(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.Content) > 4096 {
			return validation.NewInvalidInputError("Panel content must be less than 4096 characters")
		}

		return nil
	}
}

func validateChannelId(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		for _, ch := range ctx.Channels {
			if ch.Id == ctx.Data.ChannelId && (ch.Type == channel.ChannelTypeGuildText || ch.Type == channel.ChannelTypeGuildNews) {
				return nil
			}
		}

		return validation.NewInvalidInputError("Panel channel not found")
	}
}

func validateCategory(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		for _, ch := range ctx.Channels {
			if ch.Id == ctx.Data.CategoryId && ch.Type == channel.ChannelTypeGuildCategory {
				return nil
			}
		}

		return validation.NewInvalidInputError("Invalid ticket category")
	}
}

func validateEmoji(c PanelValidationContext) validation.ValidationFunc {
	return func() error {
		emoji := c.Data.Emoji

		if emoji.IsCustomEmoji {
			if emoji.Id == nil {
				return validation.NewInvalidInputError("Custom emoji was missing ID")
			}

			ctx, cancel := context.WithTimeout(context.Background(), app.DefaultTimeout)
			defer cancel()

			resolvedEmoji, err := c.BotContext.GetGuildEmoji(ctx, c.GuildId, *emoji.Id)
			if err != nil {
				return err
			}

			if resolvedEmoji.Id.Value == 0 {
				return validation.NewInvalidInputError("Emoji not found")
			}

			if resolvedEmoji.Name != emoji.Name {
				return validation.NewInvalidInputError("Emoji name mismatch")
			}
		} else {
			if len(emoji.Name) == 0 {
				return validation.NewInvalidInputError("Emoji name was empty")
			}

			// Convert from :emoji: to unicode if we need to
			name := strings.TrimSpace(emoji.Name)
			name = strings.Replace(name, ":", "", -1)

			unicode, ok := utils.GetEmoji(name)
			if !ok {
				return validation.NewInvalidInputError("Invalid emoji")
			}

			emoji.Name = unicode
		}

		return nil
	}
}

var urlRegex = regexp.MustCompile(`^https?://([-a-zA-Z0-9@:%._+~#=]{1,256})\.[a-zA-Z0-9()]{1,63}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)$`)

func validateNullableUrl(url *string) validation.ValidationFunc {
	return func() error {
		if url != nil && (len(*url) > 255 || !urlRegex.MatchString(*url)) {
			return validation.NewInvalidInputError("Invalid URL")
		}

		return nil
	}
}

func validateImageUrl(ctx PanelValidationContext) validation.ValidationFunc {
	return validateNullableUrl(ctx.Data.ImageUrl)
}

func validateThumbnailUrl(ctx PanelValidationContext) validation.ValidationFunc {
	return validateNullableUrl(ctx.Data.ThumbnailUrl)
}

func validateButtonStyle(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if ctx.Data.ButtonStyle < component.ButtonStylePrimary && ctx.Data.ButtonStyle > component.ButtonStyleDanger {
			return validation.NewInvalidInputError("Invalid button style")
		}

		return nil
	}
}

func validateButtonLabel(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if len(ctx.Data.ButtonLabel) > 80 {
			return validation.NewInvalidInputError("Button label must be less than 80 characters")
		}

		return nil
	}
}

func validatedNullableFormId(guildId uint64, formId *int) validation.ValidationFunc {
	return func() error {
		if formId == nil {
			return nil
		}

		form, ok, err := dbclient.Client.Forms.Get(context.Background(), *formId)
		if err != nil {
			return err
		}

		if !ok {
			return validation.NewInvalidInputError("Form not found")
		}

		if form.GuildId != guildId {
			return validation.NewInvalidInputError("Guild ID mismatch when validating form")
		}

		return nil
	}
}

func validateFormId(ctx PanelValidationContext) validation.ValidationFunc {
	return validatedNullableFormId(ctx.GuildId, ctx.Data.FormId)
}

// Check premium on the worker side to maintain settings if user unsubscribes and later resubscribes
func validateExitSurveyFormId(ctx PanelValidationContext) validation.ValidationFunc {
	return validatedNullableFormId(ctx.GuildId, ctx.Data.ExitSurveyFormId)
}

func validateTeams(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		// Query does not work nicely if there are no teams created in the guild, but if the user submits no teams,
		// then the input is guaranteed to be valid. Teams array excludes default team.
		if len(ctx.Data.Teams) == 0 {
			return nil
		}

		ok, err := dbclient.Client.SupportTeam.AllTeamsExistForGuild(context.Background(), ctx.GuildId, ctx.Data.Teams)
		if err != nil {
			return err
		}

		if !ok {
			return validation.NewInvalidInputError("Invalid support team")
		}

		return nil
	}
}

var placeholderPattern = regexp.MustCompile(`%(\w+)%`)

// Discord filters out illegal characters (such as +, $, ") when creating the channel for us
func validateNamingScheme(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		if ctx.Data.NamingScheme == nil {
			return nil
		}

		if len(*ctx.Data.NamingScheme) > 100 {
			return validation.NewInvalidInputError("Naming scheme must be less than 100 characters")
		}

		// Validate placeholders used
		validPlaceholders := []string{"id", "username", "nickname", "id_padded"}
		for _, match := range placeholderPattern.FindAllStringSubmatch(*ctx.Data.NamingScheme, -1) {
			if len(match) < 2 { // Infallible
				return errors.New("Infallible: Regex match length was < 2")
			}

			placeholder := match[1]
			if !utils.Contains(validPlaceholders, placeholder) {
				return validation.NewInvalidInputError(fmt.Sprintf("Invalid naming scheme placeholder: %s", placeholder))
			}
		}

		return nil
	}
}

func validateWelcomeMessage(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		wm := ctx.Data.WelcomeMessage

		if wm == nil || wm.Title != nil || wm.Description != nil || len(wm.Fields) > 0 || wm.ImageUrl != nil || wm.ThumbnailUrl != nil {
			return nil
		}

		return validation.NewInvalidInputError("Welcome message has no content")
	}
}

func validateAccessControlList(ctx PanelValidationContext) validation.ValidationFunc {
	return func() error {
		acl := ctx.Data.AccessControlList

		if len(acl) == 0 {
			return validation.NewInvalidInputError("Access control list is empty")
		}

		if len(acl) > 10 {
			return validation.NewInvalidInputError("Access control list cannot have more than 10 roles")
		}

		roles := utils.ToSet(utils.Map(ctx.Roles, utils.RoleToId))

		if roles.Size() != len(ctx.Roles) {
			return validation.NewInvalidInputError("Duplicate roles in access control list")
		}

		everyoneRoleFound := false
		for _, rule := range acl {
			if rule.RoleId == ctx.GuildId {
				everyoneRoleFound = true
			}

			if rule.Action != database.AccessControlActionDeny && rule.Action != database.AccessControlActionAllow {
				return validation.NewInvalidInputErrorf("Invalid access control action \"%s\"", rule.Action)
			}

			if !roles.Contains(rule.RoleId) {
				return validation.NewInvalidInputErrorf("Invalid role %d in access control list not found in the guild", rule.RoleId)
			}
		}

		if !everyoneRoleFound {
			return validation.NewInvalidInputError("Access control list does not contain @everyone rule")
		}

		return nil
	}
}

package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/TicketsBot/worker/bot/customisation"
	"github.com/TicketsBot/worker/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"time"
)

type (
	Settings struct {
		database.Settings
		ClaimSettings     database.ClaimSettings     `json:"claim_settings"`
		AutoCloseSettings AutoCloseData              `json:"auto_close"`
		TicketPermissions database.TicketPermissions `json:"ticket_permissions"`
		Colours           ColourMap                  `json:"colours"`

		WelcomeMessage    string                `json:"welcome_message"`
		TicketLimit       uint8                 `json:"ticket_limit"`
		Category          uint64                `json:"category,string"`
		ArchiveChannel    *uint64               `json:"archive_channel,string"`
		NamingScheme      database.NamingScheme `json:"naming_scheme"`
		UsersCanClose     bool                  `json:"users_can_close"`
		CloseConfirmation bool                  `json:"close_confirmation"`
		FeedbackEnabled   bool                  `json:"feedback_enabled"`
		Language          *string               `json:"language"`
	}

	AutoCloseData struct {
		Enabled                 bool  `json:"enabled"`
		SinceOpenWithNoResponse int64 `json:"since_open_with_no_response"`
		SinceLastMessage        int64 `json:"since_last_message"`
		OnUserLeave             bool  `json:"on_user_leave"`
	}

	ColourMap map[customisation.Colour]utils.HexColour
)

func GetSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings Settings

	group, _ := errgroup.WithContext(context.Background())

	// main settings
	group.Go(func() (err error) {
		settings.Settings, err = dbclient.Client.Settings.Get(ctx, guildId)
		return
	})

	// claim settings
	group.Go(func() (err error) {
		settings.ClaimSettings, err = dbclient.Client.ClaimSettings.Get(ctx, guildId)
		return
	})

	// auto close settings
	group.Go(func() error {
		tmp, err := dbclient.Client.AutoClose.Get(ctx, guildId)
		if err != nil {
			return err
		}

		settings.AutoCloseSettings = convertToAutoCloseData(tmp)
		return nil
	})

	// ticket permissions
	group.Go(func() (err error) {
		settings.TicketPermissions, err = dbclient.Client.TicketPermissions.Get(ctx, guildId)
		return
	})

	// colour map
	group.Go(func() (err error) {
		settings.Colours, err = getColourMap(guildId)
		return
	})

	// welcome message
	group.Go(func() (err error) {
		settings.WelcomeMessage, err = dbclient.Client.WelcomeMessages.Get(ctx, guildId)
		if err == nil && settings.WelcomeMessage == "" {
			settings.WelcomeMessage = "Thank you for contacting support.\nPlease describe your issue and await a response."
		}

		return
	})

	// ticket limit
	group.Go(func() (err error) {
		settings.TicketLimit, err = dbclient.Client.TicketLimit.Get(ctx, guildId)
		if err == nil && settings.TicketLimit == 0 {
			settings.TicketLimit = 5 // Set default
		}

		return
	})

	// category
	group.Go(func() (err error) {
		settings.Category, err = dbclient.Client.ChannelCategory.Get(ctx, guildId)
		return
	})

	// archive channel
	group.Go(func() (err error) {
		settings.ArchiveChannel, err = dbclient.Client.ArchiveChannel.Get(ctx, guildId)
		return
	})

	// allow users to close
	group.Go(func() (err error) {
		settings.UsersCanClose, err = dbclient.Client.UsersCanClose.Get(ctx, guildId)
		return
	})

	// naming scheme
	group.Go(func() (err error) {
		settings.NamingScheme, err = dbclient.Client.NamingScheme.Get(ctx, guildId)
		return
	})

	// close confirmation
	group.Go(func() (err error) {
		settings.CloseConfirmation, err = dbclient.Client.CloseConfirmation.Get(ctx, guildId)
		return
	})

	// close confirmation
	group.Go(func() (err error) {
		settings.FeedbackEnabled, err = dbclient.Client.FeedbackEnabled.Get(ctx, guildId)
		return
	})

	// language
	group.Go(func() error {
		locale, err := dbclient.Client.ActiveLanguage.Get(ctx, guildId)
		if err != nil {
			return err
		}

		if locale != "" {
			settings.Language = utils.Ptr(locale)
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// short_code -> local_name
	type MinimalLocale struct {
		IsoShortCode string `json:"iso_short_code"`
		LocalName    string `json:"local_name"`
	}

	locales := make([]MinimalLocale, len(i18n.Locales))
	for i, locale := range i18n.Locales {
		locales[i] = MinimalLocale{
			IsoShortCode: locale.IsoShortCode,
			LocalName:    locale.LocalName,
		}
	}

	ctx.JSON(200, struct {
		Settings
		Locales []MinimalLocale `json:"locales"`
	}{
		Settings: settings,
		Locales:  locales,
	})
}

func getColourMap(guildId uint64) (ColourMap, error) {
	raw, err := dbclient.Client.CustomColours.GetAll(context.Background(), guildId)
	if err != nil {
		return nil, err
	}

	colours := make(ColourMap)
	for id, hex := range raw {
		if !utils.Exists(activeColours, customisation.Colour(id)) {
			continue
		}

		colours[customisation.Colour(id)] = utils.HexColour(hex)
	}

	for _, id := range activeColours {
		if _, ok := colours[id]; !ok {
			colours[id] = utils.HexColour(customisation.DefaultColours[id])
		}
	}

	return colours, nil
}

func convertToAutoCloseData(settings database.AutoCloseSettings) (body AutoCloseData) {
	body.Enabled = settings.Enabled

	if settings.SinceOpenWithNoResponse != nil {
		body.SinceOpenWithNoResponse = int64(*settings.SinceOpenWithNoResponse / time.Second)
	}

	if settings.SinceLastMessage != nil {
		body.SinceLastMessage = int64(*settings.SinceLastMessage / time.Second)
	}

	if settings.OnUserLeave != nil {
		body.OnUserLeave = *settings.OnUserLeave
	}

	return
}

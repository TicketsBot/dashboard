package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/rest"
	"golang.org/x/sync/errgroup"
)

func UpdateSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings Settings
	if err := ctx.BindJSON(&settings); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get a list of all channel IDs
	channels := cache.Instance.GetGuildChannels(guildId)

	// TODO: Errors
	err := settings.updateSettings(guildId)
	validPrefix := settings.updatePrefix(guildId)
	validWelcomeMessage := settings.updateWelcomeMessage(guildId)
	validTicketLimit := settings.updateTicketLimit(guildId)
	validArchiveChannel := settings.updateArchiveChannel(channels, guildId)
	validCategory := settings.updateCategory(channels, guildId)
	validNamingScheme := settings.updateNamingScheme(guildId)
	settings.updatePingEveryone(guildId)
	settings.updateUsersCanClose(guildId)
	settings.updateCloseConfirmation(guildId)
	settings.updateFeedbackEnabled(guildId)

	ctx.JSON(200, gin.H{
		"prefix":          validPrefix,
		"welcome_message": validWelcomeMessage,
		"ticket_limit":    validTicketLimit,
		"archive_channel": validArchiveChannel,
		"category":        validCategory,
		"naming_scheme":   validNamingScheme,
		"error":           err,
	})
}

func (s *Settings) updateSettings(guildId uint64) error {
	if err := s.Validate(guildId); err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error {
		return dbclient.Client.Settings.Set(guildId, s.Settings)
	})

	group.Go(func() error {
		return setOpenCommandPermissions(guildId, s.DisableOpenCommand)
	})

	return group.Wait()
}

var validAutoArchive = []int{60, 1440, 4320, 10080}

func (s *Settings) Validate(guildId uint64) error {
	group, _ := errgroup.WithContext(context.Background())

	// Validate panel from same guild
	group.Go(func() error {
		if s.ContextMenuPanel != nil {
			panelId := *s.ContextMenuPanel

			panel, err := dbclient.Client.Panel.GetById(panelId)
			if err != nil {
				return err
			}

			if guildId != panel.GuildId {
				return fmt.Errorf("guild ID doesn't match")
			}
		}

		return nil
	})

	group.Go(func() error {
		valid := false
		for _, duration := range validAutoArchive {
			if duration == s.Settings.ThreadArchiveDuration {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("Invalid thread auto archive duration")
		}

		return nil
	})

	group.Go(func() error {
		if s.Settings.UseThreads {
			return fmt.Errorf("threads are disabled")
		} else {
			return nil
		}
	})

	group.Go(func() error {
        if s.Settings.OverflowCategoryId != nil {
			ch, ok := cache.Instance.GetChannel(*s.Settings.OverflowCategoryId)
			if !ok {
				return fmt.Errorf("Invalid overflow category")
			}

			if ch.GuildId != guildId {
				return fmt.Errorf("Overflow category guild ID does not match")
			}

			if ch.Type != channel.ChannelTypeGuildCategory {
				return fmt.Errorf("Overflow category is not a category")
			}
        }

		return nil
    })

	return group.Wait()
}

func setOpenCommandPermissions(guildId uint64, disabled bool) error {
	ctx, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return err
	}

	commands, err := rest.GetGlobalCommands(ctx.Token, ctx.RateLimiter, ctx.BotId)
	if err != nil {
		return err
	}

	var commandId uint64
	for _, cmd := range commands {
		if cmd.Name == "open" {
			commandId = cmd.Id
			break
		}
	}

	if commandId == 0 {
		return errors.New("open command not found")
	}

	data := rest.CommandWithPermissionsData{
		Id:            commandId,
		ApplicationId: ctx.BotId,
		GuildId:       guildId,
		Permissions: []interaction.ApplicationCommandPermissions{
			{
				Id:         guildId,
				Type:       interaction.ApplicationCommandPermissionTypeRole,
				Permission: !disabled,
			},
		},
	}

	_, err = rest.EditCommandPermissions(ctx.Token, ctx.RateLimiter, ctx.BotId, guildId, commandId, data)
	return err
}

func (s *Settings) updatePrefix(guildId uint64) bool {
	if s.Prefix == "" || len(s.Prefix) > 8 {
		return false
	}

	go dbclient.Client.Prefix.Set(guildId, s.Prefix)
	return true
}

func (s *Settings) updateWelcomeMessage(guildId uint64) bool {
	if s.WelcomeMessaage == "" || len(s.WelcomeMessaage) > 1000 {
		return false
	}

	go dbclient.Client.WelcomeMessages.Set(guildId, s.WelcomeMessaage)
	return true
}

func (s *Settings) updateTicketLimit(guildId uint64) bool {
	if s.TicketLimit > 10 || s.TicketLimit < 1 {
		return false
	}

	go dbclient.Client.TicketLimit.Set(guildId, s.TicketLimit)
	return true
}

func (s *Settings) updateCategory(channels []channel.Channel, guildId uint64) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == s.Category && ch.Type == channel.ChannelTypeGuildCategory {
			valid = true
			break
		}
	}

	if !valid {
		return false
	}

	go dbclient.Client.ChannelCategory.Set(guildId, s.Category)
	return true
}

func (s *Settings) updateArchiveChannel(channels []channel.Channel, guildId uint64) bool {
	if s.ArchiveChannel == nil {
		go dbclient.Client.ArchiveChannel.Set(guildId, nil)
		return true
	}

	var valid bool
	for _, ch := range channels {
		if ch.Id == *s.ArchiveChannel && ch.Type == channel.ChannelTypeGuildText {
			valid = true
			break
		}
	}

	if !valid {
		return false
	}

	go dbclient.Client.ArchiveChannel.Set(guildId, s.ArchiveChannel)
	return true
}

var validScheme = []database.NamingScheme{database.Id, database.Username}

func (s *Settings) updateNamingScheme(guildId uint64) bool {
	var valid bool
	for _, scheme := range validScheme {
		if scheme == s.NamingScheme {
			valid = true
			break
		}
	}

	if !valid {
		return false
	}

	go dbclient.Client.NamingScheme.Set(guildId, s.NamingScheme)
	return true
}

func (s *Settings) updatePingEveryone(guildId uint64) {
	go dbclient.Client.PingEveryone.Set(guildId, s.PingEveryone)
}

func (s *Settings) updateUsersCanClose(guildId uint64) {
	go dbclient.Client.UsersCanClose.Set(guildId, s.UsersCanClose)
}

func (s *Settings) updateCloseConfirmation(guildId uint64) {
	go dbclient.Client.CloseConfirmation.Set(guildId, s.CloseConfirmation)
}

func (s *Settings) updateFeedbackEnabled(guildId uint64) {
	go dbclient.Client.FeedbackEnabled.Set(guildId, s.FeedbackEnabled)
}

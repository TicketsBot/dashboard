package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
)

func UpdateSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings Settings
	if err := ctx.BindJSON(&settings); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	// Get a list of all channel IDs
	channels := cache.Instance.GetGuildChannels(guildId)

	// TODO: Errors
	settings.updateSettings(guildId)
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
		"prefix": validPrefix,
		"welcome_message": validWelcomeMessage,
		"ticket_limit": validTicketLimit,
		"archive_channel": validArchiveChannel,
		"category": validCategory,
		"naming_scheme": validNamingScheme,
	})
}

// TODO: Return error
func (s *Settings) updateSettings(guildId uint64) {
	go dbclient.Client.Settings.Set(guildId, s.Settings)
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
	var valid bool
	for _, ch := range channels {
		if ch.Id == s.ArchiveChannel && ch.Type == channel.ChannelTypeGuildText {
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

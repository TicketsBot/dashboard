package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
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

	// Get prefix
	validPrefix := settings.updatePrefix(guildId)
	validWelcomeMessage := settings.updateWelcomeMessage(guildId)
	validTicketLimit := settings.updateTicketLimit(guildId)
	validArchiveChannel := settings.updateArchiveChannel(channels, guildId)
	validCategory := settings.updateCategory(channels, guildId)
	validNamingScheme := settings.updateNamingScheme(guildId)
	settings.updatePingEveryone(guildId)
	settings.updateUsersCanClose(guildId)

	ctx.JSON(200, gin.H{
		"prefix": validPrefix,
		"welcome_message": validWelcomeMessage,
		"ticket_limit": validTicketLimit,
		"archive_channel": validArchiveChannel,
		"category": validCategory,
		"naming_scheme": validNamingScheme,
	})
}

func (s *Settings) updatePrefix(guildId uint64) bool {
	if s.Prefix == "" || len(s.Prefix) > 8 {
		return false
	}

	go table.UpdatePrefix(guildId, s.Prefix)
	return true
}

func (s *Settings) updateWelcomeMessage(guildId uint64) bool {
	if s.WelcomeMessaage == "" || len(s.WelcomeMessaage) > 1000 {
		return false
	}

	go table.UpdateWelcomeMessage(guildId, s.WelcomeMessaage)
	return true
}

func (s *Settings) updateTicketLimit(guildId uint64) bool {
	if s.TicketLimit > 10 || s.TicketLimit < 1 {
		return false
	}

	go table.UpdateTicketLimit(guildId, s.TicketLimit)
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

	go table.UpdateChannelCategory(guildId, s.Category)
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

	go table.UpdateArchiveChannel(guildId, s.ArchiveChannel)
	return true
}

var validScheme = []table.NamingScheme{table.Id, table.Username}
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

	go table.SetTicketNamingScheme(guildId, s.NamingScheme)
	return true
}

func (s *Settings) updatePingEveryone(guildId uint64) {
	go table.UpdatePingEveryone(guildId, s.PingEveryone)
}

func (s *Settings) updateUsersCanClose(guildId uint64) {
	go table.SetUserCanClose(guildId, s.UsersCanClose)
}

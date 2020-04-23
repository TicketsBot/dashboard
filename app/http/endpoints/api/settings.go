package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/gin-gonic/gin"
)

type Settings struct {
	Prefix          string             `json:"prefix"`
	WelcomeMessaage string             `json:"welcome_message"`
	TicketLimit     int                `json:"ticket_limit"`
	Category        uint64             `json:"category,string"`
	ArchiveChannel  uint64             `json:"archive_channel,string"`
	NamingScheme    table.NamingScheme `json:"naming_scheme"`
	PingEveryone    bool               `json:"ping_everyone"`
	UsersCanClose   bool               `json:"users_can_close"`
}

func GetSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	prefix := make(chan string)
	go table.GetPrefix(guildId, prefix)

	welcomeMessage := make(chan string)
	go table.GetWelcomeMessage(guildId, welcomeMessage)

	ticketLimit := make(chan int)
	go table.GetTicketLimit(guildId, ticketLimit)

	category := make(chan uint64)
	go table.GetChannelCategory(guildId, category)

	archiveChannel := make(chan uint64)
	go table.GetArchiveChannel(guildId, archiveChannel)

	allowUsersToClose := make(chan bool)
	go table.IsUserCanClose(guildId, allowUsersToClose)

	namingScheme := make(chan table.NamingScheme)
	go table.GetTicketNamingScheme(guildId, namingScheme)

	pingEveryone := make(chan bool)
	go table.GetPingEveryone(guildId, pingEveryone)

	ctx.JSON(200, Settings{
		Prefix:          <-prefix,
		WelcomeMessaage: <-welcomeMessage,
		TicketLimit:     <-ticketLimit,
		Category:        <-category,
		ArchiveChannel:  <-archiveChannel,
		NamingScheme:    <-namingScheme,
		PingEveryone:    <-pingEveryone,
		UsersCanClose:   <-allowUsersToClose,
	})
}

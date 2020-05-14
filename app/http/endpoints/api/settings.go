package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Settings struct {
	Prefix          string                `json:"prefix"`
	WelcomeMessaage string                `json:"welcome_message"`
	TicketLimit     uint8                 `json:"ticket_limit"`
	Category        uint64                `json:"category,string"`
	ArchiveChannel  uint64                `json:"archive_channel,string"`
	NamingScheme    database.NamingScheme `json:"naming_scheme"`
	PingEveryone    bool                  `json:"ping_everyone"`
	UsersCanClose   bool                  `json:"users_can_close"`
}

func GetSettingsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	var prefix, welcomeMessage string
	var ticketLimit uint8
	var category, archiveChannel uint64
	var allowUsersToClose, pingEveryone bool
	var namingScheme database.NamingScheme

	group, _ := errgroup.WithContext(context.Background())

	// prefix
	group.Go(func() (err error) {
		prefix, err = dbclient.Client.Prefix.Get(guildId)
		return
	})

	// welcome message
	group.Go(func() (err error) {
		welcomeMessage, err = dbclient.Client.WelcomeMessages.Get(guildId)
		return
	})

	// ticket limit
	group.Go(func() (err error) {
		ticketLimit, err = dbclient.Client.TicketLimit.Get(guildId)
		return
	})

	// category
	group.Go(func() (err error) {
		category, err = dbclient.Client.ChannelCategory.Get(guildId)
		return
	})

	// archive channel
	group.Go(func() (err error) {
		archiveChannel, err = dbclient.Client.ArchiveChannel.Get(guildId)
		return
	})

	// allow users to close
	group.Go(func() (err error) {
		allowUsersToClose, err = dbclient.Client.UsersCanClose.Get(guildId)
		return
	})

	// ping everyone
	group.Go(func() (err error) {
		pingEveryone, err = dbclient.Client.PingEveryone.Get(guildId)
		return
	})

	// naming scheme
	group.Go(func() (err error) {
		namingScheme, err = dbclient.Client.NamingScheme.Get(guildId)
		return
	})

	if err := group.Wait(); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, Settings{
		Prefix:          prefix,
		WelcomeMessaage: welcomeMessage,
		TicketLimit:     ticketLimit,
		Category:        category,
		ArchiveChannel:  archiveChannel,
		NamingScheme:    namingScheme,
		PingEveryone:    pingEveryone,
		UsersCanClose:   allowUsersToClose,
	})
}

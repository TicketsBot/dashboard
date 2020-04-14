package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"strconv"
)

func UpdateSettingsHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userId := utils.GetUserId(store)

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		guild, _ := cache.Instance.GetGuild(guildId, false)

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get CSRF token
		csrfCorrect := ctx.PostForm("csrf") == store.Get("csrf").(string)
		if !csrfCorrect {
			ctx.Redirect(302, "/")
			return
		}

		// Get prefix
		prefix := ctx.PostForm("prefix")
		prefixValid := false
		if prefix != "" && len(prefix) < 8 {
			table.UpdatePrefix(guildId, prefix)
			prefixValid = true
		}

		// Get welcome message
		welcomeMessageValid := false
		welcomeMessage := ctx.PostForm("welcomeMessage")
		if welcomeMessage != "" && len(welcomeMessage) < 1000 {
			table.UpdateWelcomeMessage(guildId, welcomeMessage)
			welcomeMessageValid = true
		}

		// Get ticket limit
		var limit int
		limitStr := ctx.PostForm("ticketlimit")

		// Verify input is an int and overwrite default limit
		if utils.IsInt(limitStr) {
			limit, _ = strconv.Atoi(limitStr)
		}

		// Update limit, or get current limit if user input is invalid
		ticketLimitValid := false
		if limitStr != "" && utils.IsInt(limitStr) && limit >= 1 && limit <= 10 {
			table.UpdateTicketLimit(guildId, limit)
			ticketLimitValid = true
		}

		// Ping everyone
		pingEveryone := ctx.PostForm("pingeveryone") == "on"
		table.UpdatePingEveryone(guildId, pingEveryone)

		// Get a list of actual category IDs
		channels := cache.Instance.GetGuildChannels(guildId)

		// Update category
		if categoryId, err := strconv.ParseUint(ctx.PostForm("category"), 10, 64); err == nil {
			for _, ch := range channels {
				if ch.Id == categoryId { // compare ID
					if ch.Type == channel.ChannelTypeGuildCategory { // verify we're dealing with a category
						table.UpdateChannelCategory(guildId, categoryId)
					}
					break
				}
			}
		}

		// Archive channel
		if archiveChannelId, err := strconv.ParseUint(ctx.PostForm("archivechannel"), 10, 64); err == nil {
			for _, ch := range channels {
				if ch.Id == archiveChannelId { // compare ID
					if ch.Type == channel.ChannelTypeGuildText { // verify channel type
						table.UpdateArchiveChannel(guildId, archiveChannelId)
					}
					break
				}
			}
		}

		// Users can close
		usersCanClose := ctx.PostForm("userscanclose") == "on"
		table.SetUserCanClose(guildId, usersCanClose)

		// Get naming scheme
		namingScheme := table.NamingScheme(ctx.PostForm("namingscheme"))
		isValidScheme := false
		for _, validNamingScheme := range table.Schemes {
			if validNamingScheme == namingScheme {
				isValidScheme = true
				break
			}
		}

		if isValidScheme {
			go table.SetTicketNamingScheme(guildId, namingScheme)
		}

		ctx.Redirect(302, fmt.Sprintf("/manage/%d/settings?validPrefix=%t&validWelcomeMessage=%t&validTicketLimit=%t", guildId, prefixValid, welcomeMessageValid, ticketLimitValid))
	} else {
		ctx.Redirect(302, "/login")
	}
}

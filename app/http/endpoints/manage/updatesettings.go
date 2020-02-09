package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UpdateSettingsHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userIdStr := store.Get("userid").(string)
		userId, err := utils.GetUserId(store)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseInt(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		var guild objects.Guild
		for _, g := range table.GetGuilds(userIdStr) {
			if g.Id == guildIdStr {
				guild = g
				break
			}
		}

		// Verify the user has permissions to be here
		if !utils.Contains(config.Conf.Admins, userIdStr) && !guild.Owner && !table.IsAdmin(guildId, userId) {
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
		categories := make(chan []table.Channel)
		go table.GetCategories(guildId, categories)

		// Convert to category IDs
		var categoryIds []string
		for _, category := range <-categories {
			categoryIds = append(categoryIds, strconv.Itoa(int(category.ChannelId)))
		}

		// Update category
		categoryStr := ctx.PostForm("category")
		if utils.Contains(categoryIds, categoryStr) {
			// Error is impossible, as we check it's a valid channel already
			category, _ := strconv.ParseInt(categoryStr, 10, 64)
			table.UpdateChannelCategory(guildId, category)
		}

		// Archive channel
		// Create a list of IDs
		var channelIds []string
		for _, c := range guild.Channels {
			channelIds = append(channelIds, c.Id)
		}

		// Update or archive channel
		archiveChannelStr := ctx.PostForm("archivechannel")
		if utils.Contains(channelIds, archiveChannelStr) {
			// Error is impossible, as we check it's a valid channel already
			parsed, _ := strconv.ParseInt(archiveChannelStr, 10, 64)
			table.UpdateArchiveChannel(guildId, parsed)
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

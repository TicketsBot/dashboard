package manage

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SettingsHandler(ctx *gin.Context) {
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
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get settings from database
		prefix := table.GetPrefix(guildId)
		welcomeMessage := table.GetWelcomeMessage(guildId)
		limit := table.GetTicketLimit(guildId)
		pingEveryone := table.GetPingEveryone(guildId)
		archiveChannel := table.GetArchiveChannel(guildId)
		categoryId := table.GetChannelCategory(guildId)

		namingSchemeChan := make(chan table.NamingScheme)
		go table.GetTicketNamingScheme(guildId, namingSchemeChan)
		namingScheme := <-namingSchemeChan

		// /users/@me/guilds doesn't return channels, so we have to get them for the specific guild
		channelsChan := make(chan []table.Channel)
		go table.GetCachedChannelsByGuild(guildId, channelsChan)
		channels := <-channelsChan

		// Get a list of actual category IDs
		categoriesChan := make(chan []table.Channel)
		go table.GetCategories(guildId, categoriesChan)
		categories := <-categoriesChan

		// Archive channel
		// Create a list of IDs
		var channelIds []string
		for _, c := range channels {
			channelIds = append(channelIds, strconv.Itoa(int(c.ChannelId)))
		}

		panelSettings := table.GetPanelSettings(guildId)

		// Users can close
		usersCanCloseChan := make(chan bool)
		go table.IsUserCanClose(guildId, usersCanCloseChan)
		usersCanClose := <-usersCanCloseChan

		invalidPrefix := ctx.Query("validPrefix") == "false"
		invalidWelcomeMessage := ctx.Query("validWelcomeMessage") == "false"
		invalidTicketLimit := ctx.Query("validTicketLimit") == "false"

		ctx.HTML(200, "manage/settings", gin.H{
			"name":           store.Get("name").(string),
			"guildId":        guildIdStr,
			"avatar": store.Get("avatar").(string),
			"prefix":         prefix,
			"welcomeMessage": welcomeMessage,
			"ticketLimit":    limit,
			"categories":     categories,
			"activecategory": categoryId,
			"channels": channels,
			"archivechannel": archiveChannel,
			"invalidPrefix": invalidPrefix,
			"invalidWelcomeMessage": invalidWelcomeMessage,
			"invalidTicketLimit": invalidTicketLimit,
			"csrf": store.Get("csrf").(string),
			"pingEveryone": pingEveryone,
			"paneltitle": panelSettings.Title,
			"panelcontent": panelSettings.Content,
			"panelcolour": strconv.FormatInt(int64(panelSettings.Colour), 16),
			"usersCanClose": usersCanClose,
			"namingScheme": string(namingScheme),
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

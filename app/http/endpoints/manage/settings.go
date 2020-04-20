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

func SettingsHandler(ctx *gin.Context) {
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

		// Check the bot is in the guild
		guild, isInGuild := cache.Instance.GetGuild(guildId, false)
		if !isInGuild {
			ctx.Redirect(302, fmt.Sprintf("https://invite.ticketsbot.net/?guild_id=%s&disable_guild_select=true&response_type=code&scope=bot%%20identify&redirect_uri=%s", guildIdStr, config.Conf.Server.BaseUrl))
			return
		}

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, userId, isAdmin)
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

		// get guild channels from cache
		channels := cache.Instance.GetGuildChannels(guildId)

		// separate out categories
		categories := make([]channel.Channel, 0)
		for _, ch := range channels {
			if ch.Type == channel.ChannelTypeGuildCategory {
				categories = append(categories, ch)
			}
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
			"name":                  store.Get("name").(string),
			"guildId":               guildIdStr,
			"avatar":                store.Get("avatar").(string),
			"prefix":                prefix,
			"welcomeMessage":        welcomeMessage,
			"ticketLimit":           limit,
			"categories":            categories,
			"activecategory":        categoryId,
			"channels":              channels,
			"archivechannel":        archiveChannel,
			"invalidPrefix":         invalidPrefix,
			"invalidWelcomeMessage": invalidWelcomeMessage,
			"invalidTicketLimit":    invalidTicketLimit,
			"csrf":                  store.Get("csrf").(string),
			"pingEveryone":          pingEveryone,
			"paneltitle":            panelSettings.Title,
			"panelcontent":          panelSettings.Content,
			"panelcolour":           strconv.FormatInt(int64(panelSettings.Colour), 16),
			"usersCanClose":         usersCanClose,
			"namingScheme":          string(namingScheme),
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

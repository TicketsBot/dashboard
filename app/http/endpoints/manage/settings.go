package manage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/template"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	guildendpoint "github.com/TicketsBot/GoPanel/utils/discord/endpoints/guild"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/apex/log"
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
		if !guild.Owner && !table.IsAdmin(guildId, userId) {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get prefix
		prefix := ctx.Query("prefix")
		if prefix == "" {
			prefix = table.GetPrefix(guildId)
		} else {
			table.UpdatePrefix(guildId, prefix)
		}

		// Get welcome message
		welcomeMessage := ctx.Query("welcomeMessage")
		if welcomeMessage == "" {
			welcomeMessage = table.GetWelcomeMessage(guildId)
		} else {
			table.UpdateWelcomeMessage(guildId, welcomeMessage)
		}

		// Get ticket limit
		limitStr := ctx.Query("ticketlimit")
		limit := 5

		// Verify input is an int and overwrite default limit
		if utils.IsInt(limitStr) {
			limit, _ = strconv.Atoi(limitStr)
		}

		// Update limit, or get current limit if user input is invalid
		if limitStr == "" || !utils.IsInt(limitStr) {
			limit = table.GetTicketLimit(guildId)
		} else {
			table.UpdateTicketLimit(guildId, limit)
		}

		// /users/@me/guilds doesn't return channels, so we have to get them for the specific guild
		if len(guild.Channels) == 0 {
			var channels []objects.Channel
			endpoint := guildendpoint.GetGuildChannels(int(guildId))
			err = endpoint.Request(store, nil, nil, &channels)

			if err != nil {
				// Not in guild
			} else {
				guild.Channels = channels

				// Update cache of categories now that we have them
				guilds := table.GetGuilds(userIdStr)

				// Get index of guild
				index := -1
				for i, g := range guilds {
					if g.Id == guild.Id {
						index = i
						break
					}
				}

				if index != -1 {
					// Delete
					guilds = append(guilds[:index], guilds[index+1:]...)

					// Insert updated guild
					guilds = utils.Insert(guilds, index, guild)

					marshalled, err := json.Marshal(guilds)
					if err != nil {
						log.Error(err.Error())
					} else {
						table.UpdateGuilds(userIdStr, base64.StdEncoding.EncodeToString(marshalled))
					}
				}
			}
		}

		// Get a list of actual category IDs
		categories := guild.GetCategories()

		// Convert to category IDs
		var categoryIds []string
		for _, c := range categories {
			categoryIds = append(categoryIds, c.Id)
		}

		categoryStr := ctx.Query("category")
		var category int64

		// Verify category ID is an int and set default category ID
		if utils.IsInt(categoryStr) {
			category, _ = strconv.ParseInt(categoryStr, 10, 64)
		}

		// Update category, or get current category if user input is invalid
		if categoryStr == "" || !utils.IsInt(categoryStr) || !utils.Contains(categoryIds, categoryStr) {
			category = table.GetChannelCategory(guildId)
		} else {
			table.UpdateChannelCategory(guildId, category)
		}

		var formattedCategories []map[string]interface{}
		for _, c := range categories {
			formattedCategories = append(formattedCategories, map[string]interface{}{
				"categoryid":   c.Id,
				"categoryname": c.Name,
				"active":  c.Id == strconv.Itoa(int(category)),
			})
		}

		// Archive channel
		// Create a list of IDs
		var channelIds []string
		for _, c := range guild.Channels {
			channelIds = append(channelIds, c.Id)
		}

		// Update or get current archive channel if blank or invalid
		var archiveChannel int64
		archiveChannelStr := ctx.Query("archivechannel")

		// Verify category ID is an int and set default category ID
		if utils.IsInt(archiveChannelStr) {
			archiveChannel, _ = strconv.ParseInt(archiveChannelStr, 10, 64)
		}

		if archiveChannelStr == "" || !utils.IsInt(archiveChannelStr) || !utils.Contains(channelIds, archiveChannelStr)  {
			archiveChannel = table.GetArchiveChannel(guildId)
		} else {
			table.UpdateArchiveChannel(guildId, archiveChannel)
		}

		// Format channels for templating
		var formattedChannels []map[string]interface{}
		for _, c := range guild.Channels {
			if c.Id == strconv.Itoa(int(archiveChannel)) {
				fmt.Println(c.Name)
			}
			if c.Type == 0 {
				formattedChannels = append(formattedChannels, map[string]interface{}{
					"channelid": c.Id,
					"channelname": c.Name,
					"active": c.Id == strconv.Itoa(int(archiveChannel)),
				})
			}
		}

		utils.Respond(ctx, template.TemplateSettings.Render(map[string]interface{}{
			"name":           store.Get("name").(string),
			"guildId":        guildIdStr,
			"prefix":         prefix,
			"welcomeMessage": welcomeMessage,
			"ticketLimit":    limit,
			"categories":     formattedCategories,
			"channels": formattedChannels,
		}))
	} else {
		ctx.Redirect(302, "/login")
	}
}

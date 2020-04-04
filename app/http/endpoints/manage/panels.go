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

type wrappedPanel struct {
	MessageId    uint64
	ChannelName  string
	Title        string
	Content      string
	CategoryName string
}

func PanelHandler(ctx *gin.Context) {
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
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
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

		// Get active panels
		panelChan := make(chan []table.Panel)
		go table.GetPanelsByGuild(guildId, panelChan)
		panels := <-panelChan

		// Get channels
		channelsChan := make(chan []table.Channel)
		go table.GetCachedChannelsByGuild(guildId, channelsChan)
		channels := <-channelsChan

		// Get default panel settings
		settings := table.GetPanelSettings(guildId)

		// Convert to wrapped panels
		wrappedPanels := make([]wrappedPanel, 0)
		for _, panel := range panels {
			wrapper := wrappedPanel{
				MessageId:    panel.MessageId,
				Title:        panel.Title,
				Content:      panel.Content,
				CategoryName: "",
			}

			if panel.Title == "" {
				wrapper.Title = settings.Title
			}
			if panel.Content == "" {
				wrapper.Content = settings.Content
			}

			// Get channel name & category name
			for _, guildChannel := range channels {
				if guildChannel.ChannelId == panel.ChannelId {
					wrapper.ChannelName = guildChannel.Name
				} else if guildChannel.ChannelId == panel.TargetCategory {
					wrapper.CategoryName = guildChannel.Name
				}
			}

			wrappedPanels = append(wrappedPanels, wrapper)
		}

		// Format channels to be text channels only
		channelMap := make(map[uint64]string)
		for _, channel := range channels {
			if channel.Type == 0 {
				channelMap[channel.ChannelId] = channel.Name
			}
		}

		// Get categories & format
		categories := make(map[uint64]string)
		for _, channel := range channels {
			if channel.Type == 4 {
				categories[channel.ChannelId] = channel.Name
			}
		}

		// Get is premium
		isPremiumChan := make(chan bool)
		go utils.IsPremiumGuild(store, guildId, isPremiumChan)
		isPremium := <-isPremiumChan

		ctx.HTML(200, "manage/panels", gin.H{
			"name":       store.Get("name").(string),
			"guildId":    guildIdStr,
			"csrf":       store.Get("csrf").(string),
			"avatar":     store.Get("avatar").(string),
			"baseUrl":    config.Conf.Server.BaseUrl,
			"panelcount": len(panels),
			"premium":    isPremium,
			"panels":     wrappedPanels,
			"channels":   channelMap,
			"categories": categories,

			"validTitle":    ctx.Query("validTitle") != "true",
			"validContent":  ctx.Query("validContent") != "false",
			"validColour":   ctx.Query("validColour") != "false",
			"validChannel":  ctx.Query("validChannel") != "false",
			"validCategory": ctx.Query("validCategory") != "false",
			"validReaction": ctx.Query("validReaction") != "false",
			"created":       ctx.Query("created") == "true",
			"metQuota":      ctx.Query("metQuota") == "true",
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

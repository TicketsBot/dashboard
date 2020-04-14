package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"strconv"
	"strings"
)

func PanelCreateHandler(ctx *gin.Context) {
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

		// Get if the guild is premium
		premiumChan := make(chan bool)
		go utils.IsPremiumGuild(store, guildId, premiumChan)
		premium := <-premiumChan

		// Check the user hasn't met their panel quota
		if !premium {
			panels := make(chan []table.Panel)
			go table.GetPanelsByGuild(guildId, panels)
			if len(<-panels) > 1 {
				ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?metQuota=true", guildId))
				return
			}
		}

		// Validate title
		title := ctx.PostForm("title")
		if len(title) == 0 || len(title) > 255 {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validTitle=false", guildId))
			return
		}

		// Validate content
		content := ctx.PostForm("content")
		if len(content) == 0 || len(content) > 1024 {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validContent=false", guildId))
			return
		}

		// Validate colour
		validColour := true
		panelColourHex := strings.Replace(ctx.PostForm("colour"), "#", "", -1)
		panelColour, err := strconv.ParseUint(panelColourHex, 16, 32)
		if err != nil {
			validColour = false
			panelColour = 0x23A31A
		}

		// Validate channel
		channelIdStr := ctx.PostForm("channel")
		channelId, err := strconv.ParseUint(channelIdStr, 10, 64); if err != nil {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validChannel=false", guildId))
			return
		}

		validChannel := make(chan bool)
		go validateChannel(guildId, channelId, validChannel)
		if !<-validChannel {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validChannel=false", guildId))
			return
		}

		// Validate category
		categoryStr := ctx.PostForm("categories")
		categoryId, err := strconv.ParseUint(categoryStr, 10, 64); if err != nil {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validCategory=false", guildId))
			return
		}

		validCategory := make(chan bool)
		go validateCategory(guildId, categoryId, validCategory)
		if !<-validCategory {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validCategory=false", guildId))
			return
		}

		// Validate reaction emote
		reaction := strings.ToLower(ctx.PostForm("reaction"))
		if len(reaction) == 0 || len(reaction) > 32 {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validReaction=false", guildId))
			return
		}
		reaction = strings.Replace(reaction, ":", "", -1)

		emoji := utils.GetEmojiByName(reaction)
		if emoji == "" {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?validReaction=false", guildId))
			return
		}

		settings := table.Panel{
			ChannelId:      channelId,
			GuildId:        guildId,
			Title:          title,
			Content:        content,
			Colour:         int(panelColour),
			TargetCategory: categoryId,
			ReactionEmote:  emoji,
		}

		go messagequeue.Client.PublishPanelCreate(settings)

		ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?created=true&validColour=%t", guildId, validColour))
	} else {
		ctx.Redirect(302, "/login")
	}
}

func validateChannel(guildId, channelId uint64, res chan bool) {
	// Compare channel IDs
	validChannel := false
	for _, guildChannel := range cache.Instance.GetGuildChannels(guildId) {
		if guildChannel.Id == channelId {
			validChannel = true
			break
		}
	}

	res <- validChannel
}

func validateCategory(guildId, categoryId uint64, res chan bool) {
	// Compare ch IDs
	validCategory := false
	for _, ch := range cache.Instance.GetGuildChannels(guildId) {
		if ch.Type == channel.ChannelTypeGuildCategory && ch.Id == categoryId {
			validCategory = true
			break
		}
	}

	res <- validCategory
}

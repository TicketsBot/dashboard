package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/cache"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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

		// Get CSRF token
		csrfCorrect := ctx.PostForm("csrf") == store.Get("csrf").(string)
		if !csrfCorrect {
			ctx.Redirect(302, "/")
			return
		}

		// Get if the guild is premium
		premiumChan := make(chan bool)
		go utils.IsPremiumGuild(store, guildIdStr, premiumChan)
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
		channelId, err := strconv.ParseInt(channelIdStr, 10, 64); if err != nil {
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
		categoryId, err := strconv.ParseInt(categoryStr, 10, 64); if err != nil {
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

		go cache.Client.PublishPanelCreate(settings)

		ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels?created=true&validColour=%t", guildId, validColour))
	} else {
		ctx.Redirect(302, "/login")
	}
}

func validateChannel(guildId, channelId int64, res chan bool) {
	// Get channels from DB
	channelsChan := make(chan []table.Channel)
	go table.GetCachedChannelsByGuild(guildId, channelsChan)
	channels := <-channelsChan

	// Compare channel IDs
	validChannel := false
	for _, guildChannel := range channels {
		if guildChannel.ChannelId == channelId {
			validChannel = true
			break
		}
	}

	res <- validChannel
}

func validateCategory(guildId, categoryId int64, res chan bool) {
	// Get channels from DB
	categoriesChan := make(chan []table.Channel)
	go table.GetCategories(guildId, categoriesChan)
	categories := <-categoriesChan

	// Compare channel IDs
	validCategory := false
	for _, category := range categories {
		if category.ChannelId == categoryId {
			validCategory = true
			break
		}
	}

	res <- validCategory
}

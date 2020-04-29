package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"strings"
)

func CreatePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	var data panel

	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	data.MessageId = 0

	// Check panel quota
	premium := make(chan bool)
	go utils.IsPremiumGuild(guildId, premium)
	if !<-premium {
		panels := make(chan []table.Panel)
		go table.GetPanelsByGuild(guildId, panels)
		if len(<-panels) > 0 {
			ctx.AbortWithStatusJSON(402, gin.H{
				"success": false,
				"error":   "You have exceeded your panel quota. Purchase premium to unlock more panels.",
			})
			return
		}
	}

	if !data.verifyTitle() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel titles must be between 1 - 255 characters in length",
		})
		return
	}

	if !data.verifyContent() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel content must be between 1 - 1024 characters in length",
		})
		return
	}

	channels := cache.Instance.GetGuildChannels(guildId)

	if !data.verifyChannel(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel",
		})
		return
	}

	if !data.verifyCategory(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel category",
		})
		return
	}

	emoji, validEmoji := data.getEmoji()
	if !validEmoji {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid emoji. Simply use the emoji's name from Discord.",
		})
		return
	}

	// TODO: Move panel create logic here
	go messagequeue.Client.PublishPanelCreate(table.Panel{
		ChannelId:      data.ChannelId,
		GuildId:        guildId,
		Title:          data.Title,
		Content:        data.Content,
		Colour:         data.Colour,
		TargetCategory: data.CategoryId,
		ReactionEmote:  emoji,
	})

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

func (p *panel) verifyTitle() bool {
	return len(p.Title) > 0 && len(p.Title) < 256
}

func (p *panel) verifyContent() bool {
	return len(p.Content) > 0 && len(p.Content) < 1025
}

func (p *panel) getEmoji() (string, bool) {
	p.Emote = strings.Replace(p.Emote, ":", "", -1)

	emoji := utils.GetEmojiByName(p.Emote)
	return emoji, emoji != ""
}

func (p *panel) verifyChannel(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.ChannelId && ch.Type == channel.ChannelTypeGuildText {
			valid = true
			break
		}
	}

	return valid
}

func (p *panel) verifyCategory(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.CategoryId && ch.Type == channel.ChannelTypeGuildCategory {
			valid = true
			break
		}
	}

	return valid
}

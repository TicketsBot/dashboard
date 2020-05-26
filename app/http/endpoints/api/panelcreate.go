package api

import (
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
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
	// TODO: Whitelabel tokens & ratelimiters
	premiumTier := rpc.PremiumClient.GetTierByGuildId(guildId, true, config.Conf.Bot.Token, ratelimit.Ratelimiter)

	if premiumTier == premium.None {
		panels, err := dbclient.Client.Panel.GetByGuild(guildId)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		if len(panels) > 0 {
			ctx.AbortWithStatusJSON(402, gin.H{
				"success": false,
				"error":   "You have exceeded your panel quota. Purchase premium to unlock more panels.",
			})
			return
		}
	}

	if !data.doValidations(ctx, guildId) {
		return
	}

	msgId, err := data.sendEmbed(premiumTier > premium.None)
	if err != nil {
		if err == request.ErrForbidden {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   "I do not have permission to send messages in the specified channel",
			})
		} else {
			// TODO: Most appropriate error?
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		return
	}

	// Add reaction
	emoji, _ := data.getEmoji() // already validated
	if err = rest.CreateReaction(config.Conf.Bot.Token, ratelimit.Ratelimiter, data.ChannelId, msgId, emoji); err != nil {
		if err == request.ErrForbidden {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   "I do not have permission to add reactions in the specified channel",
			})
		} else {
			// TODO: Most appropriate error?
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		return
	}

	// Store in DB
	panel := database.Panel{
		MessageId:      msgId,
		ChannelId:      data.ChannelId,
		GuildId:        guildId,
		Title:          data.Title,
		Content:        data.Content,
		Colour:         int32(data.Colour),
		TargetCategory: data.CategoryId,
		ReactionEmote:  emoji,
	}

	if err = dbclient.Client.Panel.Create(panel); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message_id": strconv.FormatUint(msgId, 10),
	})
}

func (p *panel) doValidations(ctx *gin.Context, guildId uint64) bool {
	if !p.verifyTitle() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel titles must be between 1 - 255 characters in length",
		})
		return false
	}

	if !p.verifyContent() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel content must be between 1 - 1024 characters in length",
		})
		return false
	}

	channels := cache.Instance.GetGuildChannels(guildId)

	if !p.verifyChannel(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel",
		})
		return false
	}

	if !p.verifyCategory(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel category",
		})
		return false
	}

	_, validEmoji := p.getEmoji()
	if !validEmoji {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid emoji. Simply use the emoji's name from Discord.",
		})
		return false
	}

	return true
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

func (p *panel) sendEmbed(isPremium bool) (messageId uint64, err error) {
	e := embed.NewEmbed().
		SetTitle(p.Title).
		SetDescription(p.Content).
		SetColor(int(p.Colour))

	if !isPremium {
		// TODO: Don't harcode
		e.SetFooter("Powered by ticketsbot.net", "https://cdn.discordapp.com/avatars/508391840525975553/ac2647ffd4025009e2aa852f719a8027.png?size=256")
	}

	var msg message.Message
	msg, err = rest.CreateMessage(config.Conf.Bot.Token, ratelimit.Ratelimiter, p.ChannelId, rest.CreateMessageData{Embed: e})
	if err != nil {
		return
	}

	messageId = msg.Id
	return
}

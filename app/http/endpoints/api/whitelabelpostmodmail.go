package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WhitelabelModmailPost(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Get bot
	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ensure bot exists
	if bot.BotId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"error": "No bot found",
		})
		return
	}

	// Parse status
	var guildId uint64
	{
		var data map[string]string
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "No guild ID provided",
			})
			return
		}

		guildIdStr, ok := data["guild"]
		if !ok {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "No guild ID provided",
			})
			return
		}

		guildId, err = strconv.ParseUint(guildIdStr, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "Invalid guild ID provided",
			})
			return
		}
	}

	if guildId == 0 {
		if err := database.Client.ModmailForcedGuilds.Delete(bot.BotId); err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
			return
		}
	}

	// verify that the bot is in the specified guild
	guilds, err := database.Client.WhitelabelGuilds.GetGuilds(bot.BotId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	var found bool
	for _, botGuild := range guilds {
		if botGuild == guildId {
			found = true
			break
		}
	}

	if !found {
		ctx.JSON(400, gin.H{
			"success": false,
			"error": "The bot isn't in the guild your provided",
		})
		return
	}

	if err := database.Client.ModmailForcedGuilds.Set(bot.BotId, guildId); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"bot": bot,
	})
}

package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WhitelabelGet(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Check if this is a different token
	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if bot.BotId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "No bot found",
		})
	} else {
		// Get status
		status, err := database.Client.WhitelabelStatuses.Get(bot.BotId)
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Get forced modmail guild
		forcedGuild, err := database.Client.ModmailForcedGuilds.Get(bot.BotId)
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"success":              true,
			"id":                   strconv.FormatUint(bot.BotId, 10),
			"status":               status,
			"modmail_forced_guild": strconv.FormatUint(forcedGuild, 10),
		})
	}
}

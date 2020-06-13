package api

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/common/statusupdates"
	"github.com/gin-gonic/gin"
)

func WhitelabelStatusPost(ctx *gin.Context) {
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
	var status string
	{
		var data map[string]string
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "No status provided",
			})
			return
		}

		var ok bool
		status, ok = data["status"]
		if !ok {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "No status provided",
			})
			return
		}

		if len(status) == 0 || len(status) > 255 {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": "Status must be between 1-255 characters in length",
			})
			return
		}
	}

	if err := database.Client.WhitelabelStatuses.Set(bot.BotId, status); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	go statusupdates.Publish(messagequeue.Client.Client, bot.BotId)

	ctx.JSON(200, gin.H{
		"success": true,
		"bot": bot,
	})
}

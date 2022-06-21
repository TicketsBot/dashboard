package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func WhitelabelGetPublicKey(ctx *gin.Context) {
	type data struct {
		PublicKey string `json:"public_key"`
	}

	userId := ctx.Keys["userid"].(uint64)

	// Get bot
	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Ensure bot exists
	if bot.BotId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "No bot found",
		})
		return
	}

	key, err := database.Client.WhitelabelKeys.Get(bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if key == "" {
		ctx.JSON(404, gin.H{
			"success": false,
		})
	} else {
		ctx.JSON(200, gin.H{
			"success": true,
			"key":     key,
		})
	}
}

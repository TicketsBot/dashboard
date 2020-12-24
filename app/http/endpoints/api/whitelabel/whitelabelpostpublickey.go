package api

import (
	"encoding/hex"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

func WhitelabelPostPublicKey(ctx *gin.Context) {
	type data struct {
		PublicKey string `json:"public_key"`
	}

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
			"error":   "No bot found",
		})
		return
	}

	// Parse status
	var body data
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "No public key provided",
		})
		return
	}

	bytes, err := hex.DecodeString(body.PublicKey)
	if err != nil || len(bytes) != 32 {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid public key",
		})
		return
	}

	if err := database.Client.WhitelabelKeys.Set(bot.BotId, body.PublicKey); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

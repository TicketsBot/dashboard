package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/common/tokenchange"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
)

func WhitelabelPost(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Get token
	var data map[string]interface{}
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Missing token",
		})
		return
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Missing token",
		})
		return
	}

	// Validate token + get bot ID
	bot, err := rest.GetCurrentUser(token, nil)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !bot.Bot {
		ctx.JSON(400, gin.H{
			"success": false,
			"error": "Token is not of a bot user",
		})
		return
	}

	// Check if this is a different token
	existing, err := dbclient.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	/*if existing.Token == token {
		// Respond with 200 to prevent information disclosure attack
		ctx.JSON(200, gin.H{
			"success": true,
			"bot": bot,
		})
		return
	}*/

	if err = dbclient.Client.Whitelabel.Set(database.WhitelabelBot{
		UserId: userId,
		BotId:  bot.Id,
		Token:  token,
	}); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	tokenChangeData := tokenchange.TokenChangeData{
		Token: token,
		NewId: bot.Id,
		OldId: existing.BotId,
	}

	if err := tokenchange.PublishTokenChange(messagequeue.Client.Client, tokenChangeData); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"bot": bot,
	})
}

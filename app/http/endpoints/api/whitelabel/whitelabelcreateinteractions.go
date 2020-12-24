package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	command "github.com/TicketsBot/worker/bot/command/impl"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"time"
)

func WhitelabelCreateInteractions(ctx *gin.Context) {
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

	botContext, err := botcontext.ContextForGuild(0)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	for _, cmd := range command.Commands {
		properties := cmd.Properties()

		if properties.MessageOnly || properties.AdminOnly || properties.HelperOnly {
			continue
		}

		option := command.BuildOption(cmd)

		data := rest.CreateCommandData{
			Name:        option.Name,
			Description: option.Description,
			Options:     option.Options,
		}

		if _, err := rest.CreateGlobalCommand(bot.Token, botContext.RateLimiter, bot.BotId, data); err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
			return
		}

		time.Sleep(time.Second)
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/worker/bot/command/impl/admin"
	"github.com/TicketsBot/worker/bot/command/manager"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/rest"
	"time"
)

// TODO: Refactor
func GetWhitelabelCreateInteractions() func(*gin.Context) {
	cm := new(manager.CommandManager)
	cm.RegisterCommands()

	return func(ctx *gin.Context) {

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

		// Cooldown
		key := fmt.Sprintf("tickets:interaction-create-cooldown:%d", bot.BotId)

		// try to set first, prevent race condition
		wasSet, err := messagequeue.Client.SetNX(key, 1, time.Minute).Result()
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// on cooldown, tell user how long left
		if !wasSet {
			expiration, err := messagequeue.Client.TTL(key).Result()
			if err != nil {
				ctx.JSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			ctx.JSON(400, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Interaction creation on cooldown, please wait another %d seconds", int64(expiration.Seconds())),
			})

			return
		}

		botContext, err := botcontext.ContextForGuild(0)
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var interactions []rest.CreateCommandData
		for _, cmd := range cm.GetCommands() {
			properties := cmd.Properties()

			if properties.MessageOnly || properties.AdminOnly || properties.HelperOnly || properties.MainBotOnly {
				continue
			}

			option := admin.BuildOption(cmd)

			data := rest.CreateCommandData{
				Name:        option.Name,
				Description: option.Description,
				Options:     option.Options,
				Type:        interaction.ApplicationCommandTypeChatInput,
			}

			interactions = append(interactions, data)
		}

		if _, err = rest.ModifyGlobalCommands(bot.Token, botContext.RateLimiter, bot.BotId, interactions); err == nil {
			ctx.JSON(200, gin.H{
				"success": true,
			})
		} else {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
}

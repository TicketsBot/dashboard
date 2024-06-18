package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/worker/bot/command/manager"
	"github.com/gin-gonic/gin"
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
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Ensure bot exists
		if bot.BotId == 0 {
			ctx.JSON(404, utils.ErrorStr("No bot found"))
			return
		}

		// Cooldown
		key := fmt.Sprintf("tickets:interaction-create-cooldown:%d", bot.BotId)

		// try to set first, prevent race condition
		wasSet, err := redis.Client.SetNX(redis.DefaultContext(), key, 1, time.Minute).Result()
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// on cooldown, tell user how long left
		if !wasSet {
			expiration, err := redis.Client.TTL(redis.DefaultContext(), key).Result()
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			ctx.JSON(400, utils.ErrorStr(fmt.Sprintf("Interaction creation on cooldown, please wait another %d seconds", int64(expiration.Seconds()))))

			return
		}

		botContext, err := botcontext.ContextForGuild(0)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		commands, _ := cm.BuildCreatePayload(true, nil)

		// TODO: Use proper context
		if _, err = rest.ModifyGlobalCommands(context.Background(), bot.Token, botContext.RateLimiter, bot.BotId, commands); err == nil {
			ctx.JSON(200, utils.SuccessResponse)
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}
	}
}

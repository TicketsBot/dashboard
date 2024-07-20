package api

import (
	"context"
	"errors"
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
		bot, err := database.Client.Whitelabel.GetByUserId(ctx, userId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Ensure bot exists
		if bot.BotId == 0 {
			ctx.JSON(404, utils.ErrorStr("No bot found"))
			return
		}

		if err := createInteractions(cm, bot.BotId, bot.Token); err != nil {
			if errors.Is(err, ErrInteractionCreateCooldown) {
				ctx.JSON(400, utils.ErrorJson(err))
			} else {
				ctx.JSON(500, utils.ErrorJson(err))
			}

			return
		}

		ctx.JSON(200, utils.SuccessResponse)
	}
}

var ErrInteractionCreateCooldown = errors.New("Interaction creation on cooldown")

func createInteractions(cm *manager.CommandManager, botId uint64, token string) error {
	// Cooldown
	key := fmt.Sprintf("tickets:interaction-create-cooldown:%d", botId)

	// try to set first, prevent race condition
	wasSet, err := redis.Client.SetNX(redis.DefaultContext(), key, 1, time.Minute).Result()
	if err != nil {
		return err
	}

	// on cooldown, tell user how long left
	if !wasSet {
		expiration, err := redis.Client.TTL(redis.DefaultContext(), key).Result()
		if err != nil {
			return err
		}

		return fmt.Errorf("%w, please wait another %d seconds", ErrInteractionCreateCooldown, int64(expiration.Seconds()))
	}

	botContext, err := botcontext.ContextForGuild(0)
	if err != nil {
		return err
	}

	commands, _ := cm.BuildCreatePayload(true, nil)

	// TODO: Use proper context
	_, err = rest.ModifyGlobalCommands(context.Background(), token, botContext.RateLimiter, botId, commands)
	return err
}

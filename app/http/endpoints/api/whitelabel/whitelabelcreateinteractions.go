package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/worker/bot/command/manager"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"net/http"
	"time"
)

// TODO: Refactor
func GetWhitelabelCreateInteractions() func(*gin.Context) {
	cm := new(manager.CommandManager)
	cm.RegisterCommands()

	return func(c *gin.Context) {
		userId := c.Keys["userid"].(uint64)

		// Get bot
		bot, err := database.Client.Whitelabel.GetByUserId(c, userId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			return
		}

		// Ensure bot exists
		if bot.BotId == 0 {
			c.JSON(404, utils.ErrorStr("No bot found"))
			return
		}

		if err := createInteractions(cm, bot.BotId, bot.Token); err != nil {
			if errors.Is(err, ErrInteractionCreateCooldown) {
				c.JSON(http.StatusTooManyRequests, utils.ErrorJson(err))
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
			}

			return
		}

		c.JSON(200, utils.SuccessResponse)
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

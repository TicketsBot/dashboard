package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/redis"
	wrapper "github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/gin-gonic/gin"
	"time"
)

func ReloadGuildsHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	key := fmt.Sprintf("tickets:dashboard:guildreload:%d", userId)
	res, err := redis.Client.SetNX(wrapper.DefaultContext(), key, 1, time.Second*10).Result()
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !res {
		ttl, err := redis.Client.TTL(wrapper.DefaultContext(), key).Result()
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// handle redis error codes
		if ttl < 0 {
			ttl = 0
		}

		ctx.JSON(429, utils.ErrorStr("You're doing this too quickly: try again in %d seconds", int(ttl.Seconds())))
		return
	}

	store, err := session.Store.Get(userId)
	if err != nil {
		if err == session.ErrNoSession {
			ctx.JSON(401, gin.H{
				"success": false,
				"auth":    true,
			})
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	// What does this do?
	if store.Expiry > time.Now().Unix() {
		res, err := discord.RefreshToken(store.RefreshToken)
		if err != nil { // Tell client to re-authenticate
			ctx.JSON(200, gin.H{
				"success":                 false,
				"reauthenticate_required": true,
			})
			return
		}

		store.AccessToken = res.AccessToken
		store.RefreshToken = res.RefreshToken
		store.Expiry = time.Now().Unix() + int64(res.ExpiresIn)

		if err := session.Store.Set(userId, store); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	if err := utils.LoadGuilds(store.AccessToken, userId); err != nil {
		// TODO: Log to sentry
		// Tell client to reauth, needs a 200 or client will display error
		ctx.JSON(200, gin.H{
			"success":                 false,
			"reauthenticate_required": true,
		})
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

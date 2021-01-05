package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

func ReloadGuildsHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	key := fmt.Sprintf("tickets:dashboard:guildreload:%d", userId)
	res, err := messagequeue.Client.SetNX(key, 1, time.Second*10).Result()
	if err != nil {
		ctx.JSON(500, utils.ErrorToResponse(err))
		return
	}

	if !res {
		ttl, err := messagequeue.Client.TTL(key).Result()
		if err != nil {
			ctx.JSON(500, utils.ErrorToResponse(err))
			return
		}

		// handle redis error codes
		if ttl < 0 {
			ttl = 0
		}

		ctx.JSON(429, utils.ErrorStr("You're doing this too quickly: try again in %d seconds", int(ttl.Seconds())))
		return
	}

	store := sessions.Default(ctx)
	if store == nil {
		ctx.JSON(200, gin.H{
			"success": false,
			"reauthenticate_required": true,
		})
		return
	}

	accessToken := store.Get("access_token").(string)
	expiry := store.Get("expiry").(int64)
	if expiry > (time.Now().UnixNano() / int64(time.Second)) {
		res, err := discord.RefreshToken(store.Get("refresh_token").(string))
		if err != nil { // Tell client to re-authenticate
			ctx.JSON(200, gin.H{
				"success": false,
				"reauthenticate_required": true,
			})
			return
		}

		accessToken = res.AccessToken

		store.Set("access_token", res.AccessToken)
		store.Set("refresh_token", res.RefreshToken)
		store.Set("expiry", (time.Now().UnixNano()/int64(time.Second))+int64(res.ExpiresIn))
		store.Save()
	}

	if err := utils.LoadGuilds(accessToken, userId); err != nil {
		ctx.JSON(500, utils.ErrorToResponse(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

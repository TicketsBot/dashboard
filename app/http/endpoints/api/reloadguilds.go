package api

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/redis"
	wrapper "github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"net/http"
	"time"
)

func ReloadGuildsHandler(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	key := fmt.Sprintf("tickets:dashboard:guildreload:%d", userId)
	res, err := redis.Client.SetNX(wrapper.DefaultContext(), key, 1, time.Second*10).Result()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !res {
		ttl, err := redis.Client.TTL(wrapper.DefaultContext(), key).Result()
		if err != nil {
			c.JSON(500, utils.ErrorJson(err))
			return
		}

		// handle redis error codes
		if ttl < 0 {
			ttl = 0
		}

		c.JSON(429, utils.ErrorStr("You're doing this too quickly: try again in %d seconds", int(ttl.Seconds())))
		return
	}

	store, err := session.Store.Get(userId)
	if err != nil {
		if err == session.ErrNoSession {
			c.JSON(401, gin.H{
				"success": false,
				"auth":    true,
			})
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		return
	}

	// What does this do?
	if store.Expiry > time.Now().Unix() {
		res, err := rest.RefreshToken(c, nil, config.Conf.Oauth.Id, config.Conf.Oauth.Secret, store.RefreshToken)
		if err != nil { // Tell client to re-authenticate
			c.JSON(200, gin.H{
				"success":                 false,
				"reauthenticate_required": true,
			})
			return
		}

		store.AccessToken = res.AccessToken
		store.RefreshToken = res.RefreshToken
		store.Expiry = time.Now().Unix() + int64(res.ExpiresIn)

		if err := session.Store.Set(userId, store); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	guilds, err := utils.LoadGuilds(c, store.AccessToken, userId)
	if err != nil {
		var oauthError request.OAuthError
		if errors.As(err, &oauthError) {
			if oauthError.ErrorCode == "invalid_grant" {
				// Tell client to reauth, needs a 200 or client will display error
				c.JSON(http.StatusOK, gin.H{
					"success":                 false,
					"reauthenticate_required": true,
				})
				return
			}
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"guilds":  guilds,
	})
}

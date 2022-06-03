package api

import (
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
)

func SessionHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	store, err := session.Store.Get(userId)
	if err != nil {
		if err == session.ErrNoSession {
			ctx.JSON(404, gin.H{
				"success": false,
				"error":   err.Error(),
				"auth":    true,
			})
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	tier, err := rpc.PremiumClient.GetTierByUser(userId, false)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"username":   store.Name,
		"avatar":     store.Avatar,
		"whitelabel": tier >= premium.Whitelabel,
	})
}

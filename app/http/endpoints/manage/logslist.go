package manage

import (
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func LogsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	store, err := session.Store.Get(userId)
	if err != nil {
		if err == session.ErrNoSession {
			ctx.JSON(401, gin.H{
				"success": false,
				"auth": true,
			})
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	ctx.HTML(200, "manage/logs", gin.H{
		"name":         store.Name,
		"guildId":      guildId,
		"avatar":       store.Avatar,
		"baseUrl":      config.Conf.Server.BaseUrl,
	})
}

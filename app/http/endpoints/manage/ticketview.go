package manage

import (
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func TicketViewHandler(ctx *gin.Context) {
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

	guildId := ctx.Keys["guildid"].(uint64)

	ctx.HTML(200, "manage/ticketview", gin.H{
		"name":     store.Name,
		"guildId":  guildId,
		"avatar":   store.Avatar,
		"baseUrl":  config.Conf.Server.BaseUrl,
		"ticketId": ctx.Param("ticketId"),
	})
}

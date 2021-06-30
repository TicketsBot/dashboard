package root

import (
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	if err := session.Store.Clear(userId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

package middleware

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthenticateCookie(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		ctx.Redirect(302, "/login")
		ctx.Abort()
		return
	}

	if !utils.IsLoggedIn(store) {
		ctx.Redirect(302, "/login")
		ctx.Abort()
		return
	}

	ctx.Keys["userid"] = utils.GetUserId(store)
}

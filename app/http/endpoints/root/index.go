package root

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)

	ctx.HTML(200, "main/index", gin.H{
		"name":    store.Get("name").(string),
		"baseurl": config.Conf.Server.BaseUrl,
		"avatar":  store.Get("avatar").(string),
	})
}

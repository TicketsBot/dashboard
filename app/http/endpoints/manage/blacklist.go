package manage

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func BlacklistHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	guildId := ctx.Keys["guildid"].(uint64)

	ctx.HTML(200, "manage/blacklist", gin.H{
		"name":         store.Get("name").(string),
		"guildId":      guildId,
		"avatar":       store.Get("avatar").(string),
		"baseUrl":      config.Conf.Server.BaseUrl,
	})
}

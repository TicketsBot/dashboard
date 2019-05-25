package endpoints

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/url"
)

func LoginHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	redirect := url.QueryEscape(fmt.Sprintf("%s/callback", config.Conf.Server.BaseUrl))
	ctx.Redirect(302, fmt.Sprintf("https://discordapp.com/oauth2/authorize?response_type=code&redirect_uri=%s&scope=identify+guilds&client_id=%d", redirect, config.Conf.Oauth.Id))
}

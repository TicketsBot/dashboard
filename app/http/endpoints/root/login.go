package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
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

	if utils.IsLoggedIn(store) {
		ctx.Redirect(302, config.Conf.Server.BaseUrl)
	} else {
		redirect := url.QueryEscape(config.Conf.Oauth.RedirectUri)

		var guildsScope string
		if _, noGuilds := ctx.GetQuery("noguilds"); !noGuilds {
			guildsScope = "+guilds"
		}

		ctx.Redirect(302, fmt.Sprintf("https://discordapp.com/oauth2/authorize?response_type=code&redirect_uri=%s&scope=identify%s&client_id=%d&state=%s", redirect, guildsScope, config.Conf.Oauth.Id, ctx.Query("state")))
	}
}

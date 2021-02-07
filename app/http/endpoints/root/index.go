package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/url"
)

func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)

	if _, hasGuilds := store.Get("has_guilds").(bool); !hasGuilds {
		redirect := url.QueryEscape(config.Conf.Oauth.RedirectUri)
		ctx.Redirect(302, fmt.Sprintf("https://discordapp.com/oauth2/authorize?response_type=code&redirect_uri=%s&scope=identify+guilds&client_id=%d&state=%s", redirect, config.Conf.Oauth.Id, ctx.Query("state")))
		return
	}

	ctx.HTML(200, "main/index", gin.H{
		"name":         store.Get("name").(string),
		"baseurl":      config.Conf.Server.BaseUrl,
		"avatar":       store.Get("avatar").(string),
		"referralShow": config.Conf.Referral.Show,
		"referralLink": config.Conf.Referral.Link,
	})
}

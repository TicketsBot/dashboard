package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"net/url"
)

func IndexHandler(ctx *gin.Context) {
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

	if !store.HasGuilds {
		redirect := url.QueryEscape(config.Conf.Oauth.RedirectUri)
		ctx.Redirect(302, fmt.Sprintf("https://discordapp.com/oauth2/authorize?response_type=code&redirect_uri=%s&scope=identify+guilds&client_id=%d&state=%s", redirect, config.Conf.Oauth.Id, ctx.Query("state")))
		return
	}

	ctx.HTML(200, "main/index", gin.H{
		"name":         store.Name,
		"baseurl":      config.Conf.Server.BaseUrl,
		"avatar":       store.Avatar,
		"referralShow": config.Conf.Referral.Show,
		"referralLink": config.Conf.Referral.Link,
	})
}

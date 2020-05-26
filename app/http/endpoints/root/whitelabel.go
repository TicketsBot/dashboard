package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func WhitelabelHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	userId := store.Get("userid").(uint64)

	premiumTier := rpc.PremiumClient.GetTierByUser(userId, false)
	if premiumTier < premium.Whitelabel {
		ctx.Redirect(302, fmt.Sprintf("%s/premium", config.Conf.Server.MainSite))
		return
	}

	ctx.HTML(200, "main/whitelabel", gin.H{
		"name":    store.Get("name").(string),
		"baseurl": config.Conf.Server.BaseUrl,
		"avatar":  store.Get("avatar").(string),
	})
}

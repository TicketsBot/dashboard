package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
)

func WhitelabelHandler(ctx *gin.Context) {
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

	premiumTier := rpc.PremiumClient.GetTierByUser(userId, false)
	if premiumTier < premium.Whitelabel {
		var isForced bool
		for _, forced := range config.Conf.ForceWhitelabel {
			if forced == userId {
				isForced = true
				break
			}
		}

		if !isForced {
			ctx.Redirect(302, fmt.Sprintf("%s/premium", config.Conf.Server.MainSite))
			return
		}
	}

	ctx.HTML(200, "main/whitelabel", gin.H{
		"name":    store.Name,
		"baseurl": config.Conf.Server.BaseUrl,
		"avatar":  store.Avatar,
		"referralShow": config.Conf.Referral.Show,
		"referralLink": config.Conf.Referral.Link,
	})
}

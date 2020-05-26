package api

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
)

func PremiumHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// TODO: Whitelabel tokens & ratelimiters
	premiumTier := rpc.PremiumClient.GetTierByGuildId(guildId, true, config.Conf.Bot.Token, ratelimit.Ratelimiter)

	ctx.JSON(200, gin.H{
		"premium": premiumTier >= premium.Premium,
	})
}

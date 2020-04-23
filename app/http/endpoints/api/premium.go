package api

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func PremiumHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	isPremium := make(chan bool)
	go utils.IsPremiumGuild(guildId, isPremium)

	ctx.JSON(200, gin.H{
		"premium": <-isPremium,
	})
}

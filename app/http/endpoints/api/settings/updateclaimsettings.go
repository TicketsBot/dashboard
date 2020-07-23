package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

func PostClaimSettings(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings database.ClaimSettings
	if err := ctx.BindJSON(&settings); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	if settings.SupportCanType && !settings.SupportCanView {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": "Must be able to view channel to type",
		})
		return
	}

	if err := dbclient.Client.ClaimSettings.Set(guildId, settings); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

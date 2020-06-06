package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

func PostAutoClose(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var settings database.AutoCloseSettings
	if err := ctx.BindJSON(&settings); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	if settings.Enabled && (settings.SinceLastMessage == nil || settings.SinceOpenWithNoResponse == nil || settings.OnUserLeave == nil) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": "No time period provided",
		})
		return
	}

	if (settings.SinceOpenWithNoResponse != nil && *settings.SinceOpenWithNoResponse < 0) || (settings.SinceLastMessage != nil && *settings.SinceLastMessage < 0) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": "Negative time period provided",
		})
		return
	}

	if !settings.Enabled {
		settings.SinceLastMessage = nil
		settings.SinceOpenWithNoResponse = nil
		settings.OnUserLeave = nil
	}

	if err := dbclient.Client.AutoClose.Set(guildId, settings); err != nil {
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

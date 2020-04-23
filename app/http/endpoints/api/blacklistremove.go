package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/gin-gonic/gin"
	"strconv"
)

func RemoveBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	userId, err := strconv.ParseUint(ctx.Param("user"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	go table.RemoveBlacklist(guildId, userId)

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

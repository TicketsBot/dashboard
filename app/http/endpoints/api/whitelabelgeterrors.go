package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

func WhitelabelGetErrors(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	errors, err := database.Client.WhitelabelErrors.GetRecent(userId, 10)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"errors": errors,
	})
}

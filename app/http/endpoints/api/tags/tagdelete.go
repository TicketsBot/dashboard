package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

func DeleteTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	tagId := ctx.Param("tag")

	if tagId == "" || len(tagId) > 16 {
		ctx.JSON(400, gin.H{
			"success": false,
			"error": "Invalid tag",
		})
		return
	}

	if err := database.Client.Tag.Delete(guildId, tagId); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"success": true,
		})
	}
}

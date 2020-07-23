package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

// TODO: Make client take new structure
func TagsListHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	tags, err := database.Client.Tag.GetByGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, tags)
}

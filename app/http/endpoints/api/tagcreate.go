package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

type tag struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func CreateTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	var data tag

	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !data.verifyIdLength() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Tag ID must be 1 - 16 characters in length",
		})
		return
	}

	if !data.verifyContentLength() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Tag content must be 1 - 2000 characters in length",
		})
		return
	}

	if err := database.Client.Tag.Set(guildId, data.Id, data.Content); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"success": true,
		})
	}
}

func (t *tag) verifyIdLength() bool {
	return len(t.Id) > 0 && len(t.Id) <= 16
}

func (t *tag) verifyContentLength() bool {
	return len(t.Content) > 0 && len(t.Content) <= 2000
}

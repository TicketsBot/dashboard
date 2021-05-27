package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

func DeletePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	panel, err := database.Client.Panel.GetById(panelId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	// verify panel belongs to guild
	if panel.GuildId != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error": "Guild ID doesn't match",
		})
		return
	}

	if err :=  database.Client.Panel.Delete(panelId); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	if err := rest.DeleteMessage(botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.MessageId); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

package api

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

func DeletePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	messageId, err := strconv.ParseUint(ctx.Param("message"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	panel, err := database.Client.Panel.Get(messageId)
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

	if err :=  database.Client.Panel.Delete(messageId); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	if err := rest.DeleteMessage(config.Conf.Bot.Token, ratelimit.Ratelimiter, panel.ChannelId, panel.MessageId); err != nil {
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

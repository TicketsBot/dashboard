package api

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
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

	// verify panel belongs to guild
	panelChan := make(chan table.Panel)
	go table.GetPanel(messageId, panelChan)
	panel := <-panelChan

	if panel.GuildId != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error": "Guild ID doesn't match",
		})
		return
	}

	go table.DeletePanel(messageId)
	go rest.DeleteMessage(config.Conf.Bot.Token, ratelimit.Ratelimiter, panel.ChannelId, panel.MessageId)

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

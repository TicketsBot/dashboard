package api

import (
	"encoding/json"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
)

func ChannelsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	channels := cache.Instance.GetGuildChannels(guildId)
	encoded, err := json.Marshal(channels)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": true,
			"error": err.Error(),
		})
		return
	}

	ctx.Data(200, gin.MIMEJSON, encoded)
}

package api

import (
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
)

func ChannelsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	channels := cache.Instance.GetGuildChannels(guildId)
	if channels == nil {
		channels = make([]channel.Channel, 0) // don't serve null
	}

	ctx.JSON(200, channels)
}

package api

import (
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"sort"
)

func ChannelsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	channels := cache.Instance.GetGuildChannels(guildId)
	if channels == nil {
		channels = make([]channel.Channel, 0) // don't serve null
	} else {
		sort.Slice(channels, func(i,j int) bool {
			return channels[i].Position < channels[j].Position
		})
	}

	ctx.JSON(200, channels)
}

package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	cache2 "github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/channel"
	"sort"
)

func ChannelsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// TODO: Use proper context
	channels, err := cache.Instance.GetGuildChannels(context.Background(), guildId)
	if err != nil {
		if errors.Is(err, cache2.ErrNotFound) {
			ctx.JSON(200, make([]channel.Channel, 0))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	filtered := make([]channel.Channel, 0, len(channels))
	for _, ch := range channels {
		// Filter out threads
		if ch.Type == channel.ChannelTypeGuildNewsThread ||
			ch.Type == channel.ChannelTypeGuildPrivateThread ||
			ch.Type == channel.ChannelTypeGuildPublicThread {
			continue
		}

		filtered = append(filtered, ch)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Position < filtered[j].Position
	})

	ctx.JSON(200, filtered)
}

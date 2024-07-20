package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	cache2 "github.com/rxdn/gdl/cache"
	"strconv"
)

func WhitelabelGetGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	bot, err := database.Client.Whitelabel.GetByUserId(ctx, userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// id -> name
	guilds := make(map[string]string)
	if bot.BotId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"guilds":  guilds,
		})
		return
	}

	ids, err := database.Client.WhitelabelGuilds.GetGuilds(ctx, bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	for _, id := range ids {
		// get guild name
		// TODO: Use proper context
		guild, err := cache.Instance.GetGuild(context.Background(), id)
		if err != nil {
			if errors.Is(err, cache2.ErrNotFound) {
				continue
			} else {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		}

		guilds[strconv.FormatUint(id, 10)] = guild.Name
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"guilds":  guilds,
	})
}

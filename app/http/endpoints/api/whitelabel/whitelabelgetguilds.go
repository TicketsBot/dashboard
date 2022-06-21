package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func WhitelabelGetGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// id -> name
	guilds := make(map[string]string, 0)
	if bot.BotId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"guilds":  guilds,
		})
		return
	}

	ids, err := database.Client.WhitelabelGuilds.GetGuilds(bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	for _, id := range ids {
		// get guild name
		if guild, found := cache.Instance.GetGuild(id, false); found {
			guilds[strconv.FormatUint(id, 10)] = guild.Name
		}
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"guilds":  guilds,
	})
}

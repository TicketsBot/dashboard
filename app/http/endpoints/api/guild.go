package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func GuildHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	guild, err := botContext.GetGuild(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"id":   guild.Id,
		"name": guild.Name,
		"icon": guild.Icon,
	})
}

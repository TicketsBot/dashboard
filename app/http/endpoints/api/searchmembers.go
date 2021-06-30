package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/member"
)

func SearchMembers(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botCtx, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	query := ctx.Query("query")
	if len(query) > 32 {
		ctx.JSON(400, utils.ErrorStr("Invalid query"))
		return
	}

	var members []member.Member
	if query == "" {
		members, err = botCtx.ListMembers(guildId)
	} else {
		members, err = botCtx.SearchMembers(guildId, query)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, members)
}

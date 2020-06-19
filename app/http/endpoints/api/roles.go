package api

import (
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
)

func RolesHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	roles := cache.Instance.GetGuildRoles(guildId)

	ctx.JSON(200, gin.H{
		"success": true,
		"roles": roles,
	})
}

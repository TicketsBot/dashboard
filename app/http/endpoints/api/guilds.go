package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
)

func GetGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	userGuilds := table.GetGuilds(userId)
	adminGuilds := make([]guild.Guild, 0)
	for _, g := range userGuilds {
		fakeGuild := guild.Guild{
			Id:          g.Id,
			OwnerId:     g.OwnerId,
			Permissions: g.Permissions,
		}

		isAdmin := make(chan bool)
		go utils.IsAdmin(fakeGuild, userId, isAdmin)
		if <-isAdmin {
			adminGuilds = append(adminGuilds, g)
		}
	}

	ctx.JSON(200, adminGuilds)
}

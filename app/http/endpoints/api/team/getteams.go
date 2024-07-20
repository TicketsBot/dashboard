package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

func GetTeams(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	teams, err := dbclient.Client.SupportTeam.Get(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// prevent serving null
	if teams == nil {
		teams = make([]database.SupportTeam, 0)
	}

	ctx.JSON(200, teams)
}

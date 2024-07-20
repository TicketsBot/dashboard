package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

func CreateTeam(ctx *gin.Context) {
	type body struct {
		Name string `json:"name"`
	}

	guildId := ctx.Keys["guildid"].(uint64)

	var data body
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if len(data.Name) == 0 || len(data.Name) > 32 {
		ctx.JSON(400, utils.ErrorStr("Team name must be between 1 and 32 characters"))
		return
	}

	id, err := dbclient.Client.SupportTeam.Create(ctx, guildId, data.Name)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, database.SupportTeam{
		Id:      id,
		GuildId: guildId,
		Name:    data.Name,
	})
}

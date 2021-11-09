package customisation

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/worker/bot/customisation"
	"github.com/gin-gonic/gin"
)

// GetColours TODO: Don't depend on worker
func GetColours(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// TODO: Don't duplicate
	raw, err := dbclient.Client.CustomColours.GetAll(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	colours := make(map[customisation.Colour]utils.HexColour)
	for id, hex := range raw {
		colours[customisation.Colour(id)] = utils.HexColour(hex)
	}

	for id, hex := range customisation.DefaultColours {
		if _, ok := colours[id]; !ok {
			colours[id] = utils.HexColour(hex)
		}
	}

	ctx.JSON(200, colours)
}

package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

type deleteBody struct {
	TagId string `json:"tag_id"`
}

func DeleteTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var body deleteBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// Increase max length for characters from other alphabets
	if body.TagId == "" || len(body.TagId) > 100 {
		ctx.JSON(400, utils.ErrorStr("Invalid tag"))
		return
	}

	if err := database.Client.Tag.Delete(guildId, body.TagId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

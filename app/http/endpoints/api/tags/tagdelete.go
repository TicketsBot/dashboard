package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

func DeleteTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	type Body struct {
		TagId string `json:"tag_id"`
	}

	var body Body
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if body.TagId == "" || len(body.TagId) > 16 {
		ctx.JSON(400, utils.ErrorStr("Invalid tag"))
		return
	}

	if err := database.Client.Tag.Delete(guildId, body.TagId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
	} else {
		ctx.JSON(200, utils.SuccessResponse)
	}
}

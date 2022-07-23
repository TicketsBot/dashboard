package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

type tag struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func CreateTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// Max of 200 tags
	count, err := database.Client.Tag.GetTagCount(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if count >= 200 {
		ctx.JSON(400, utils.ErrorStr("Tag limit (200) reached"))
		return
	}

	var data tag
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if !data.verifyIdLength() {
		ctx.JSON(400, utils.ErrorStr("Tag ID must be 1 - 16 characters in length"))
		return
	}

	if !data.verifyContentLength() {
		ctx.JSON(400, utils.ErrorStr("Tag content must be 1 - 2000 characters in length"))
		return
	}

	if err := database.Client.Tag.Set(guildId, data.Id, data.Content); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
	} else {
		ctx.JSON(200, utils.SuccessResponse)
	}
}

func (t *tag) verifyIdLength() bool {
	return len(t.Id) > 0 && len(t.Id) <= 16
}

func (t *tag) verifyContentLength() bool {
	return len(t.Content) > 0 && len(t.Content) <= 2000
}

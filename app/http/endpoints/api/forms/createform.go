package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

type createFormBody struct {
	Title string `json:"title"`
}

func CreateForm(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var data createFormBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if len(data.Title) > 255 {
		ctx.JSON(400, utils.ErrorStr("Title is too long"))
        return
	}

	// 26^50 chance of collision
	customId := utils.RandString(50)

	id, err := dbclient.Client.Forms.Create(guildId, data.Title, customId)
	if err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	form := database.Form{
		Id:       id,
		GuildId:  guildId,
		Title:    data.Title,
		CustomId: customId,
	}

	ctx.JSON(200, form)
}

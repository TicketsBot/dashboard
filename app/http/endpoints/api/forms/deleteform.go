package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func DeleteForm(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	formId, err := strconv.Atoi(ctx.Param("form_id"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid form ID"))
		return
	}

	form, ok, err := dbclient.Client.Forms.Get(formId)
	if err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	if !ok {
		ctx.JSON(404, utils.ErrorStr("Form not found"))
		return
	}

	if form.GuildId != guildId {
		ctx.JSON(403, utils.ErrorStr("Form does not belong to this guild"))
        return
	}

	if err := dbclient.Client.Forms.Delete(formId); err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	ctx.JSON(200, utils.SuccessResponse)
}

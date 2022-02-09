package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func DeleteInput(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	formId, err := strconv.Atoi(ctx.Param("form_id"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid form ID"))
		return
	}

	inputId, err := strconv.Atoi(ctx.Param("input_id"))
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

	input, ok, err := dbclient.Client.FormInput.Get(inputId)
	if err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	if !ok {
		ctx.JSON(404, utils.ErrorStr("Input not found"))
        return
	}

	if input.FormId != formId {
        ctx.JSON(403, utils.ErrorStr("Input does not belong to this form"))
        return
    }

	if err := dbclient.Client.FormInput.Delete(input.Id, input.FormId); err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	ctx.JSON(200, utils.SuccessResponse)
}

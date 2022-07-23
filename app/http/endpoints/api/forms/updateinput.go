package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UpdateInput(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var data inputCreateBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if !data.Validate(ctx) {
		return
	}

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

	newInput := database.FormInput{
		Id:          inputId,
		FormId:      formId,
		CustomId:    input.CustomId,
		Style:       uint8(data.Style),
		Label:       data.Label,
		Placeholder: data.Placeholder,
		Required:    data.Required,
	}

	if err := dbclient.Client.FormInput.Update(newInput); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, newInput)
}

package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
)

type embeddedForm struct {
	database.Form
	Inputs []database.FormInput `json:"inputs"`
}

func GetForms(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	forms, err := dbclient.Client.Forms.GetForms(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	inputs, err := dbclient.Client.FormInput.GetInputsForGuild(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data := make([]embeddedForm, len(forms))
	for i, form := range forms {
		formInputs, ok := inputs[form.Id]
		if !ok {
			formInputs = make([]database.FormInput, 0)
		}

		data[i] = embeddedForm{
			Form:   form,
			Inputs: formInputs,
		}
	}

	ctx.JSON(200, data)
}

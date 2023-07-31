package forms

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/interaction/component"
	"strconv"
)

func CreateInput(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var data inputCreateBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// Validate body
	if !data.Validate(ctx) {
		return
	}

	// Parse form ID from URL
	formId, err := strconv.Atoi(ctx.Param("form_id"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid form ID"))
		return
	}

	// Get form and validate it belongs to the guild
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

	// Check there are not more than 25 inputs already
	// TODO: This is vulnerable to a race condition
	inputCount, err := getFormInputCount(formId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if inputCount >= 5 {
		ctx.JSON(400, utils.ErrorStr("A form cannot have more than 5 inputs"))
		return
	}

	// 2^30 chance of collision
	customId, err := utils.RandString(30)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	formInputId, err := dbclient.Client.FormInput.Create(formId, customId, uint8(data.Style), data.Label, data.Placeholder, data.Required)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, database.FormInput{
		Id:          formInputId,
		FormId:      formId,
		CustomId:    customId,
		Style:       uint8(data.Style),
		Label:       data.Label,
		Placeholder: data.Placeholder,
		Required:    data.Required,
	})
}

func (b *inputCreateBody) Validate(ctx *gin.Context) bool {
	if b.Style != component.TextStyleShort && b.Style != component.TextStyleParagraph {
		ctx.JSON(400, utils.ErrorStr("Invalid style"))
		return false
	}

	if len(b.Label) == 0 || len(b.Label) > 45 {
		ctx.JSON(400, utils.ErrorStr("The input label must be between 1 and 45 characters"))
		return false
	}

	if b.Placeholder != nil && len(*b.Placeholder) == 0 {
		b.Placeholder = nil
	}

	if b.Placeholder != nil && len(*b.Placeholder) > 100 {
		ctx.JSON(400, utils.ErrorStr("The placeholder cannot be more than 100 characters"))
		return false
	}

	return true
}

// TODO: Use select count()
func getFormInputCount(formId int) (int, error) {
	inputs, err := dbclient.Client.FormInput.GetInputs(formId)
	if err != nil {
		return 0, err
	}

	return len(inputs), nil
}

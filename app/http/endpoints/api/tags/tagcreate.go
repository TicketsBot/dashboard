package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

type tag struct {
	Id              string             `json:"id" validate:"required,min=1,max=16"`
	UseGuildCommand bool               `json:"use_guild_command"` // Not yet implemented
	Content         *string            `json:"content" validate:"omitempty,min=1,max=4096"`
	UseEmbed        bool               `json:"use_embed"`
	Embed           *types.CustomEmbed `json:"embed" validate:"omitempty,dive"`
}

var (
	validate          = validator.New()
	slashCommandRegex = regexp.MustCompile(`^[-_a-zA-Z0-9]{1,32}$`)
)

func CreateTag(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	// Max of 200 tags
	count, err := dbclient.Client.Tag.GetTagCount(guildId)
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

	if !data.UseEmbed {
		data.Embed = nil
	}

	// TODO: Limit command amount
	if err := validate.Struct(data); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(500, utils.ErrorStr("An error occurred while validating the integration"))
			return
		}

		formatted := "Your input contained the following errors:\n" + utils.FormatValidationErrors(validationErrors)
		ctx.JSON(400, utils.ErrorStr(formatted))
		return
	}

	if !data.verifyContent() {
		ctx.JSON(400, utils.ErrorStr("You have not provided any content for the tag"))
		return
	}

	var embed *database.CustomEmbedWithFields
	if data.Embed != nil {
		customEmbed, fields := data.Embed.IntoDatabaseStruct()
		embed = &database.CustomEmbedWithFields{
			CustomEmbed: customEmbed,
			Fields:      fields,
		}
	}

	wrapped := database.Tag{
		Id:              data.Id,
		GuildId:         guildId,
		UseGuildCommand: data.UseGuildCommand,
		Content:         data.Content,
		Embed:           embed,
	}

	if err := dbclient.Client.Tag.Set(wrapped); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

func (t *tag) verifyId() bool {
	if len(t.Id) == 0 || len(t.Id) > 16 || strings.Contains(t.Id, " ") {
		return false
	}

	if t.UseGuildCommand {
		return slashCommandRegex.MatchString(t.Id)
	} else {
		return true
	}
}

func (t *tag) verifyContent() bool {
	if t.Content != nil { // validator ensures that if this is not nil, > 0 length
		return true
	}

	if t.Embed != nil {
		if t.Embed.Description != nil || len(t.Embed.Fields) > 0 || t.Embed.ImageUrl != nil || t.Embed.ThumbnailUrl != nil {
			return true
		}
	}

	return false
}

package types

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel/embed"
)

type CustomEmbed struct {
	Title        *string        `json:"title" validate:"omitempty,min=1,max=255"`
	Description  *string        `json:"description" validate:"omitempty,min=1,max=4096"`
	Url          *string        `json:"url" validate:"omitempty,url,max=255"`
	Colour       Colour         `json:"colour" validate:"required,gte=0,lte=16777215"`
	Author       Author         `json:"author" validate:"dive"`
	ImageUrl     *string        `json:"image_url" validate:"omitempty,url,max=255"`
	ThumbnailUrl *string        `json:"thumbnail_url" validate:"omitempty,url,max=255"`
	Footer       Footer         `json:"footer" validate:"dive"`
	Timestamp    *DateTimeLocal `json:"timestamp" validate:"omitempty"`
	Fields       []Field        `json:"fields" validate:"dive,max=25"`
}

type Author struct {
	Name    *string `json:"name" validate:"omitempty,min=1,max=255"`
	IconUrl *string `json:"icon_url" validate:"omitempty,url,max=255"`
	Url     *string `json:"url" validate:"omitempty,url,max=255"`
}

type Footer struct {
	Text    *string `json:"text" validate:"omitempty,min=1,max=2048"`
	IconUrl *string `json:"icon_url" validate:"omitempty,url,max=255"`
}

type Field struct {
	Name   string `json:"name" validate:"min=1,max=255"`
	Value  string `json:"value" validate:"min=1,max=1024"`
	Inline bool   `json:"inline"`
}

func NewCustomEmbed(c *database.CustomEmbed, fields []database.EmbedField) *CustomEmbed {
	wrappedFields := make([]Field, len(fields))
	for i, field := range fields {
		wrappedFields[i] = Field{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}
	}

	return &CustomEmbed{
		Title:       c.Title,
		Description: c.Description,
		Url:         c.Url,
		Colour:      Colour(c.Colour),
		Author: Author{
			Name:    c.AuthorName,
			IconUrl: c.AuthorIconUrl,
			Url:     c.AuthorUrl,
		},
		ImageUrl:     c.ImageUrl,
		ThumbnailUrl: c.ThumbnailUrl,
		Footer: Footer{
			Text:    c.FooterText,
			IconUrl: c.FooterIconUrl,
		},
		Timestamp: NewDateTimeLocalFromPtr(c.Timestamp),
		Fields:    wrappedFields,
	}
}

func (c *CustomEmbed) IntoDatabaseStruct() (*database.CustomEmbed, []database.EmbedField) {
	fields := make([]database.EmbedField, len(c.Fields))
	for i, field := range c.Fields {
		fields[i] = database.EmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		}
	}

	return &database.CustomEmbed{
		Title:         c.Title,
		Description:   c.Description,
		Url:           c.Url,
		Colour:        uint32(c.Colour),
		AuthorName:    c.Author.Name,
		AuthorIconUrl: c.Author.IconUrl,
		AuthorUrl:     c.Author.Url,
		ImageUrl:      c.ImageUrl,
		ThumbnailUrl:  c.ThumbnailUrl,
		FooterText:    c.Footer.Text,
		FooterIconUrl: c.Footer.IconUrl,
		Timestamp:     TimeOrNil(c.Timestamp),
	}, fields
}

func (c *CustomEmbed) IntoDiscordEmbed() *embed.Embed {
	e := &embed.Embed{
		Title:       utils.ValueOrZero(c.Title),
		Description: utils.ValueOrZero(c.Description),
		Url:         utils.ValueOrZero(c.Url),
		Timestamp:   TimeOrNil(c.Timestamp),
		Color:       int(c.Colour),
	}

	if c.Footer.Text != nil {
		e.Footer = &embed.EmbedFooter{
			Text:    *c.Footer.Text,
			IconUrl: utils.ValueOrZero(c.Footer.IconUrl),
		}
	}

	if c.ImageUrl != nil {
		e.SetImage(*c.ImageUrl)
	}

	if c.ThumbnailUrl != nil {
		e.SetThumbnail(*c.ThumbnailUrl)
	}

	if c.Author.Name != nil {
		e.Author = &embed.EmbedAuthor{
			Name:    *c.Author.Name,
			IconUrl: utils.ValueOrZero(c.Author.IconUrl),
			Url:     utils.ValueOrZero(c.Author.Url),
		}
	}

	if len(c.Fields) > 0 {
		e.Fields = make([]*embed.EmbedField, len(c.Fields))

		for i, field := range c.Fields {
			e.Fields[i] = &embed.EmbedField{
				Name:   field.Name,
				Value:  field.Value,
				Inline: field.Inline,
			}
		}
	}

	return e
}

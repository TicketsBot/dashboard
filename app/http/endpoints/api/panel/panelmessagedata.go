package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
)

type panelMessageData struct {
	ChannelId uint64

	Title, Content, CustomId      string
	Colour                        int
	ImageUrl, ThumbnailUrl, Emoji *string
	ButtonStyle                   component.ButtonStyle
	ButtonLabel                   string
	IsPremium                     bool
}

func panelIntoMessageData(panel database.Panel, isPremium bool) panelMessageData {
	var emoji *string
	if panel.ReactionEmote != "" {
		emoji = &panel.ReactionEmote
	}

	return panelMessageData{
		ChannelId:    panel.ChannelId,
		Title:        panel.Title,
		Content:      panel.Content,
		CustomId:     panel.CustomId,
		Colour:       int(panel.Colour),
		ImageUrl:     panel.ImageUrl,
		ThumbnailUrl: panel.ThumbnailUrl,
		Emoji:        emoji,
		ButtonStyle:  component.ButtonStyle(panel.ButtonStyle),
		ButtonLabel:  panel.ButtonLabel,
		IsPremium:    isPremium,
	}
}

func (p *panelMessageData) send(ctx *botcontext.BotContext) (uint64, error) {
	e := embed.NewEmbed().
		SetTitle(p.Title).
		SetDescription(p.Content).
		SetColor(p.Colour)

	if p.ImageUrl != nil {
		e.SetImage(*p.ImageUrl)
	}

	if p.ThumbnailUrl != nil {
		e.SetThumbnail(*p.ThumbnailUrl)
	}

	if !p.IsPremium {
		e.SetFooter("Powered by ticketsbot.net", "https://ticketsbot.net/assets/img/logo.png")
	}

	var buttonEmoji *emoji.Emoji
	if p.Emoji != nil {
		buttonEmoji = &emoji.Emoji{
			Name: *p.Emoji,
		}
	}

	data := rest.CreateMessageData{
		Embeds: []*embed.Embed{e},
		Components: []component.Component{
			component.BuildActionRow(component.BuildButton(component.Button{
				Label:    p.ButtonLabel,
				CustomId: p.CustomId,
				Style:    p.ButtonStyle,
				Emoji:    buttonEmoji,
				Url:      nil,
				Disabled: false,
			})),
		},
	}

	msg, err := rest.CreateMessage(ctx.Token, ctx.RateLimiter, p.ChannelId, data)
	if err != nil {
		return 0, err
	}

	return msg.Id, nil
}

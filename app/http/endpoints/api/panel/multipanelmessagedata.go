package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
	"math"
)

type multiPanelMessageData struct {
	ChannelId uint64

	Title     string
	Content   string
	Colour    int
	IsPremium bool
}

func multiPanelIntoMessageData(panel database.MultiPanel, isPremium bool) multiPanelMessageData {
	return multiPanelMessageData{
		ChannelId: panel.ChannelId,
		Title:     panel.Title,
		Content:   panel.Content,
		Colour:    panel.Colour,
		IsPremium: isPremium,
	}
}

func (d *multiPanelMessageData) send(ctx *botcontext.BotContext, panels []database.Panel) (uint64, error) {
	e := embed.NewEmbed().
		SetTitle(d.Title).
		SetDescription(d.Content).
		SetColor(d.Colour)

	if !d.IsPremium {
		// TODO: Don't harcode
		e.SetFooter("Powered by ticketsbot.net", "https://ticketsbot.net/assets/img/logo.png")
	}

	buttons := make([]component.Component, len(panels))
	for i, panel := range panels {
		var buttonEmoji *emoji.Emoji
		if panel.ReactionEmote != "" {
			buttonEmoji = &emoji.Emoji{
				Name: panel.ReactionEmote,
			}
		}

		buttons[i] = component.BuildButton(component.Button{
			Label:    panel.Title,
			CustomId: panel.CustomId,
			Style:    component.ButtonStyle(panel.ButtonStyle),
			Emoji:    buttonEmoji,
		})
	}

	var rows []component.Component
	for i := 0; i <= int(math.Ceil(float64(len(buttons)/5))); i++ {
		lb := i * 5
		ub := lb + 5

		if ub >= len(buttons) {
			ub = len(buttons)
		}

		if lb >= ub {
			break
		}

		row := component.BuildActionRow(buttons[lb:ub]...)
		rows = append(rows, row)
	}

	data := rest.CreateMessageData{
		Embeds:     []*embed.Embed{e},
		Components: rows,
	}

	msg, err := rest.CreateMessage(ctx.Token, ctx.RateLimiter, d.ChannelId, data)
	if err != nil {
		return 0, err
	}

	return msg.Id, nil
}

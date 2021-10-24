package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/utils"
	"math"
)

type multiPanelMessageData struct {
	ChannelId uint64

	Title      string
	Content    string
	Colour     int
	SelectMenu bool
	IsPremium  bool
}

func multiPanelIntoMessageData(panel database.MultiPanel, isPremium bool) multiPanelMessageData {
	return multiPanelMessageData{
		ChannelId:  panel.ChannelId,
		Title:      panel.Title,
		Content:    panel.Content,
		Colour:     panel.Colour,
		SelectMenu: panel.SelectMenu,
		IsPremium:  isPremium,
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

	var components []component.Component
	if d.SelectMenu {
		options := make([]component.SelectOption, len(panels))
		for i, panel := range panels {
			var emote *emoji.Emoji
			if panel.ReactionEmote != "" {
				emote = &emoji.Emoji{
					Name: panel.ReactionEmote,
				}
			}

			options[i] = component.SelectOption{
				Label: panel.Title,
				Value: panel.CustomId,
				Emoji: emote,
			}
		}

		components = []component.Component{
			component.BuildActionRow(
				component.BuildSelectMenu(
					component.SelectMenu{
						CustomId:    "multipanel",
						Options:     options,
						Placeholder: "Select a topic...",
						MinValues:   utils.IntPtr(1),
						MaxValues:   utils.IntPtr(1),
						Disabled:    false,
					}),
			),
		}
	} else {
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

		components = rows
	}

	data := rest.CreateMessageData{
		Embeds:     []*embed.Embed{e},
		Components: components,
	}

	msg, err := rest.CreateMessage(ctx.Token, ctx.RateLimiter, d.ChannelId, data)
	if err != nil {
		return 0, err
	}

	return msg.Id, nil
}

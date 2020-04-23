package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/gin-gonic/gin"
)

type panel struct {
	ChannelId  uint64 `json:"channel_id,string"`
	MessageId  uint64 `json:"message_id,string"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Colour     uint32 `json:"colour"`
	CategoryId uint64 `json:"category_id,string"`
	Emote      string `json:"emote"`
}

func ListPanels(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	panelsChan := make(chan []table.Panel)
	go table.GetPanelsByGuild(guildId, panelsChan)
	panels := <-panelsChan

	wrapped := make([]panel, len(panels))

	for i, p := range panels {
		wrapped[i] = panel{
			ChannelId:  p.ChannelId,
			MessageId:  p.MessageId,
			Title:      p.Title,
			Content:    p.Content,
			Colour:     p.Colour,
			CategoryId: p.TargetCategory,
			Emote:      p.ReactionEmote,
		}
	}

	ctx.JSON(200, wrapped)
}

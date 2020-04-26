package api

import (
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/gin-gonic/gin"
)

func TagsListHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	wrapped := make([]tag, 0)

	tags := make(chan []table.Tag)
	go table.GetTags(guildId, tags)
	for _, t := range <-tags {
		wrapped = append(wrapped, tag{
			Id:      t.Id,
			Content: t.Content,
		})
	}

	ctx.JSON(200, wrapped)
}

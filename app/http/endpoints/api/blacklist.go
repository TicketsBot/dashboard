package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"strconv"
)

type userData struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

func GetBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	data := make(map[string]userData)

	blacklistedUsers := table.GetBlacklistNodes(guildId)
	for _, row := range blacklistedUsers {
		formattedId := strconv.FormatUint(row.User, 10)
		user, _ := cache.Instance.GetUser(row.User)

		data[formattedId] = userData{
			Username:      user.Username,
			Discriminator: fmt.Sprintf("%04d", user.Discriminator),
		}
	}

	ctx.JSON(200, data)
}

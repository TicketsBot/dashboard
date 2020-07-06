package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetPermissionLevel(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	guilds := strings.Split(ctx.Query("guilds"), ",")
	if len(guilds) > 100 {
		ctx.JSON(400, gin.H{
			"success": false,
			"error": "too many guilds",
		})
		return
	}

	// TODO: Check whether the bot is in the guild to prevent us getting maliciously 429'd

	levels := make(map[string]permission.PermissionLevel)

	for _, raw := range guilds {
		guildId, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": fmt.Sprintf("invalid guild id: %s", raw),
			})
			return
		}

		level := utils.GetPermissionLevel(guildId, userId)
		levels[strconv.FormatUint(guildId, 10)] = level
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"levels": levels,
	})
}

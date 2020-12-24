package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

	// TODO: This is insanely inefficient

	levels := make(map[string]permission.PermissionLevel)

	group, _ := errgroup.WithContext(context.Background())
	for _, raw := range guilds {
		guildId, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error": fmt.Sprintf("invalid guild id: %s", raw),
			})
			return
		}

		group.Go(func() error {
			level, err := utils.GetPermissionLevel(guildId, userId)
			levels[strconv.FormatUint(guildId, 10)] = level
			return err
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorToResponse(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"levels": levels,
	})
}

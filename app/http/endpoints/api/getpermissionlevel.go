package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPermissionLevel(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	guildId, err := strconv.ParseUint(ctx.Query("guild"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid guild ID"))
		return
	}

	// TODO: Use proper context
	permissionLevel, err := utils.GetPermissionLevel(context.Background(), guildId, userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success":          true,
		"permission_level": permissionLevel,
	})
}

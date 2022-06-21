package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UserHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	userId, err := strconv.ParseUint(ctx.Param("user"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid user ID"))
		return
	}

	var username string
	if err := cache.Instance.QueryRow(context.Background(), `SELECT "data"->>'Username' FROM users WHERE users.user_id=$1 AND EXISTS(SELECT FROM members WHERE members.guild_id=$2);`, userId, guildId).Scan(&username); err != nil {
		ctx.JSON(404, utils.ErrorStr("Not found"))
		return
	}

	ctx.JSON(200, gin.H{
		"user_id":  userId,
		"guild_id": guildId,
		"username": username,
	})
}

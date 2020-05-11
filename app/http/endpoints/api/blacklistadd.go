package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"strconv"
)

func AddBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var data userData
	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	parsedDiscrim, err := strconv.ParseInt(data.Discriminator, 10, 16)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	var targetId uint64
	if err := cache.Instance.QueryRow(context.Background(), `select users.user_id from "users" where LOWER(users.data->>'Username')=LOWER($1) AND users.data->>'Discriminator'=$2;`, data.Username, strconv.FormatInt(parsedDiscrim, 10)).Scan(&targetId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"error": "user not found",
			})
		} else {
			fmt.Println(err.Error())
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
		}
		return
	}

	// TODO: Don't blacklist staff or guild owner
	if err = database.Client.Blacklist.Add(guildId, targetId); err == nil {
		ctx.JSON(200, gin.H{
			"success": true,
			"user_id": strconv.FormatUint(targetId, 10),
		})
	} else {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
	}
}

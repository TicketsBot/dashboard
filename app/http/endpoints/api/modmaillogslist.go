package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
)

var modmailPathRegex = regexp.MustCompile(`(\d+)\/modmail\/(?:free-)?([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

type wrappedModLog struct {
	Uuid    string `json:"uuid"`
	GuildId uint64 `json:"guild_id,string"`
	UserId  uint64 `json:"user_id,string"`
}

func GetModmailLogs(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	// filter
	var userId uint64

	if userIdRaw, filterByUserId := ctx.GetQuery("userid"); filterByUserId {
		userId, err = strconv.ParseUint(userIdRaw, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"error": "Invalid user ID",
			})
			return
		}
	} else if username, filterByUsername := ctx.GetQuery("username"); filterByUsername {
		if err := cache.Instance.QueryRow(context.Background(), `select users.user_id from users where LOWER("data"->>'Username') LIKE LOWER($1) and exists(SELECT FROM members where members.guild_id=$2);`, fmt.Sprintf("%%%s%%", username), guildId).Scan(&userId); err != nil {
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"error": "User not found",
			})
			return
		}
	}

	shouldFilter := userId > 0

	start := pageLimit * (page - 1)
	end := start + pageLimit - 1

	wrapped := make([]wrappedModLog, 0)

	var archives []table.ModMailArchive
	if shouldFilter {
		archivesCh := make(chan []table.ModMailArchive)
		go table.GetModmailArchivesByUser(userId, guildId, archivesCh)
		archives = <-archivesCh
	} else {
		archivesCh := make(chan []table.ModMailArchive)
		go table.GetModmailArchivesByGuild(guildId, archivesCh)
		archives = <-archivesCh
	}

	for i := start; i < end && i < len(archives); i++ {
		wrapped = append(wrapped, wrappedModLog{
			Uuid:    archives[i].Uuid,
			GuildId: archives[i].Guild,
			UserId:  archives[i].User,
		})
	}

	ctx.JSON(200, wrapped)
}

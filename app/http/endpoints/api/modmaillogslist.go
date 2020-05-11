package api

import (
	"context"
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"regexp"
	"strconv"
)

var modmailPathRegex = regexp.MustCompile(`(\d+)\/modmail\/(?:free-)?([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

type wrappedModLog struct {
	Uuid    string `json:"uuid"`
	GuildId uint64 `json:"guild_id,string"`
	UserId  uint64 `json:"user_id,string"`
}

// TODO: Take after param
func GetModmailLogs(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	after, err := uuid.FromString(ctx.Query("after"))
	if err != nil {
		after = uuid.Nil
	}

	before, err := uuid.FromString(ctx.Query("before"))
	if err != nil {
		before = uuid.Nil
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

	wrapped := make([]wrappedModLog, 0)

	var archives []database.ModmailArchive
	if shouldFilter {
		archives, err = dbclient.Client.ModmailArchive.GetByMember(guildId, userId, pageLimit, after, before)
	} else {
		archives, err = dbclient.Client.ModmailArchive.GetByGuild(guildId, pageLimit, after, before)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	for _, archive := range archives {
		wrapped = append(wrapped, wrappedModLog{
			Uuid:    archive.Uuid.String(),
			GuildId: archive.GuildId,
			UserId:  archive.UserId,
		})
	}

	ctx.JSON(200, wrapped)
}

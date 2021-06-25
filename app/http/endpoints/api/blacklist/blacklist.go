package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
	"golang.org/x/sync/errgroup"
)

type userData struct {
	UserId        uint64             `json:"id,string"`
	Username      string             `json:"username"`
	Discriminator user.Discriminator `json:"discriminator"`
}

// TODO: Paginate
func GetBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	blacklistedUsers, err := database.Client.Blacklist.GetBlacklistedUsers(guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	data := make([]userData, len(blacklistedUsers))

	group, _ := errgroup.WithContext(context.Background())
	for i, userId := range blacklistedUsers {
		i := i
		userId := userId

		// TODO: Mass lookup
		group.Go(func() error {
			userData := userData{
				UserId: userId,
			}

			user, ok := cache.Instance.GetUser(userId)
			if ok {
				userData.Username = user.Username
				userData.Discriminator = user.Discriminator
			}

			data[i] = userData
			return nil
		})
	}

	_ = group.Wait()

	ctx.JSON(200, data)
}

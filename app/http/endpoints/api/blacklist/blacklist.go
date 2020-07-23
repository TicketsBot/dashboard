package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"strconv"
	"sync"
)

type userData struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

func GetBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	blacklistedUsers, err := database.Client.Blacklist.GetBlacklistedUsers(guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	data := make(map[string]userData)
	var lock sync.Mutex

	group, _ := errgroup.WithContext(context.Background())
	for _, userId := range blacklistedUsers {
		group.Go(func() error {
			user, _ := cache.Instance.GetUser(userId)

			lock.Lock()

			// JS cant do big ints
			data[strconv.FormatUint(userId, 10)] = userData{
				Username:      user.Username,
				Discriminator: fmt.Sprintf("%04d", user.Discriminator),
			}

			lock.Unlock()

			return nil
		})
	}

	_ = group.Wait()

	ctx.JSON(200, data)
}

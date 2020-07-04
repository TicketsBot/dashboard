package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

type wrappedGuild struct {
	Id   uint64 `json:"id,string"`
	Name string `json:"name"`
}

func GetGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	guilds, err := database.Client.UserGuilds.Get(userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	group, _ := errgroup.WithContext(context.Background())

	adminGuilds := make([]wrappedGuild, 0)
	var lock sync.Mutex

	for _, g := range guilds {
		g := g

		group.Go(func() error {
			fakeGuild := guild.Guild{
				Id:          g.GuildId,
				Owner:       g.Owner,
				Permissions: int(g.UserPermissions),
			}

			if g.Owner {
				fakeGuild.OwnerId = userId
			}

			if utils.GetPermissionLevel(g.GuildId, userId) >= permission.Support {
				lock.Lock()
				adminGuilds = append(adminGuilds, wrappedGuild{
					Id:   g.GuildId,
					Name: g.Name,
				})
				lock.Unlock()
			}
			return nil
		})
	}

	// not possible anyway but eh
	if err := group.Wait(); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// sort guilds
	sort.Slice(adminGuilds, func(i, j int) bool {
		return adminGuilds[i].Name < adminGuilds[j].Name
	})

	ctx.JSON(200, adminGuilds)
}

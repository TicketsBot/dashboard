package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

type wrappedGuild struct {
	Id   uint64 `json:"id,string"`
	Name string `json:"name"`
	Icon string `json:"icon"`
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
			// verify bot is in guild
			_, ok := cache.Instance.GetGuild(g.GuildId, false)
			if !ok {
				return nil
			}
			
			fakeGuild := guild.Guild{
				Id:          g.GuildId,
				Owner:       g.Owner,
				Permissions: uint64(g.UserPermissions),
			}

			if g.Owner {
				fakeGuild.OwnerId = userId
			}

			permLevel, err := utils.GetPermissionLevel(g.GuildId, userId)
			if err != nil {
				// If a Discord error occurs, just skip the server
				if _, ok := err.(*request.RestError); !ok {
					return err
				}
			}

			if permLevel >= permission.Support {
				lock.Lock()
				adminGuilds = append(adminGuilds, wrappedGuild{
					Id:   g.GuildId,
					Name: g.Name,
					Icon: g.Icon,
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

/*func getAdminGuilds(userId uint64) ([]uint64, error) {
	var guilds []uint64

	// get guilds owned by user
	query := `SELECT "guild_id" FROM guilds WHERE "data"->'owner_id' = '$1';`
	rows, err := cache.Instance.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var guildId uint64
		if err := rows.Scan(&guildId); err != nil {
			return nil, err
		}

		guilds = append(guilds, guildId)
	}

	database.Client.Permissions.GetSupport()
}*/

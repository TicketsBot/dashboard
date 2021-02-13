package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	syncutils "github.com/TicketsBot/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"sort"
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

	wg := syncutils.NewChannelWaitGroup()
	wg.Add(len(guilds))

	group, _ := errgroup.WithContext(context.Background())
	ch := make(chan wrappedGuild)
	for _, g := range guilds {
		g := g

		group.Go(func() error {
			defer wg.Done()

			// verify bot is in guild
			if err := cache.Instance.QueryRow(context.Background(), `SELECT 1 from guilds WHERE "guild_id" = $1`, g.GuildId).Scan(nil); err != nil {
				if err == pgx.ErrNoRows {
					return nil
				} else {
					return err
				}
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
				if _, ok := err.(request.RestError); !ok {
					return err
				}
			}

			if permLevel >= permission.Support {
				wrapped := wrappedGuild{
					Id:   g.GuildId,
					Name: g.Name,
					Icon: g.Icon,
				}

				ch <- wrapped
			}

			return nil
		})
	}

	adminGuilds := make([]wrappedGuild, 0)
	group.Go(func() error {
	loop:
		for {
			select {
			case <-wg.Wait():
				break loop
			case guild := <-ch:
				adminGuilds = append(adminGuilds, guild)
			}
		}
		return nil
	})

	_ = group.Wait() // error not possible

	// sort
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

package api

import (
	"context"
	"errors"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	syncutils "github.com/TicketsBot/common/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"sort"
)

type wrappedGuild struct {
	Id              uint64                     `json:"id,string"`
	Name            string                     `json:"name"`
	Icon            string                     `json:"icon"`
	PermissionLevel permission.PermissionLevel `json:"permission_level"`
}

func GetGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Get all guilds the user is in
	guilds, err := dbclient.Client.UserGuilds.Get(userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Get the subset of guilds that the user is in that the bot is also in
	guildIds := make([]uint64, len(guilds))
	guildMap := make(map[uint64]database.UserGuild) // Make a map of all guilds for O(1) access
	for i, guild := range guilds {
		guildIds[i] = guild.GuildId
		guildMap[guild.GuildId] = guild
	}

	botGuilds, err := getExistingGuilds(guildIds)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	wg := syncutils.NewChannelWaitGroup()
	wg.Add(len(botGuilds))

	group, _ := errgroup.WithContext(context.Background())
	ch := make(chan wrappedGuild)
	for _, guildId := range botGuilds {
		guildId := guildId
		g := guildMap[guildId]

		group.Go(func() error {
			defer wg.Done()

			// Determine the user's permission level in this guild
			var permLevel permission.PermissionLevel
			if g.Owner {
				permLevel = permission.Admin
			} else {
				permLevel, err = utils.GetPermissionLevel(context.Background(), g.GuildId, userId)
				if err != nil {
					// If a Discord error occurs, just skip the server
					var restError request.RestError
					if !errors.As(err, &restError) {
						return err
					}
				}
			}

			if permLevel >= permission.Support {
				wrapped := wrappedGuild{
					Id:              g.GuildId,
					Name:            g.Name,
					Icon:            g.Icon,
					PermissionLevel: permLevel,
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

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// sort
	sort.Slice(adminGuilds, func(i, j int) bool {
		return adminGuilds[i].Name < adminGuilds[j].Name
	})

	ctx.JSON(200, adminGuilds)
}

func getExistingGuilds(userGuilds []uint64) ([]uint64, error) {
	query := `SELECT "guild_id" from guilds WHERE "guild_id" = ANY($1);`

	userGuildsArray := &pgtype.Int8Array{}
	if err := userGuildsArray.Set(userGuilds); err != nil {
		return nil, err
	}

	rows, err := cache.Instance.Query(context.Background(), query, userGuildsArray)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var existingGuilds []uint64
	for rows.Next() {
		var guildId uint64
		if err := rows.Scan(&guildId); err != nil {
			return nil, err
		}

		existingGuilds = append(existingGuilds, guildId)
	}

	return existingGuilds, nil
}

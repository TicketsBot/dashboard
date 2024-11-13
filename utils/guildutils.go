package utils

import (
	"cmp"
	"context"
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/common/collections"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/database"
	"github.com/jackc/pgtype"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest"
	errgroup "golang.org/x/sync/errgroup"
	"slices"
	"sync"
)

type GuildDto struct {
	Id              uint64                     `json:"id,string"`
	Name            string                     `json:"name"`
	Icon            string                     `json:"icon"`
	PermissionLevel permission.PermissionLevel `json:"permission_level"`
}

func LoadGuilds(ctx context.Context, accessToken string, userId uint64) ([]GuildDto, error) {
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	data := rest.CurrentUserGuildsData{
		Limit: 200,
	}

	guilds, err := rest.GetCurrentUserGuilds(ctx, authHeader, nil, data)
	if err != nil {
		return nil, err
	}

	if err := storeGuildsInDb(ctx, userId, guilds); err != nil {
		return nil, err
	}

	userGuilds, err := getGuildIntersection(ctx, userId, guilds)
	if err != nil {
		return nil, err
	}

	group, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex
	dtos := make([]GuildDto, 0, len(userGuilds))
	for _, guild := range userGuilds {
		guild := guild

		group.Go(func() error {
			permLevel, err := GetPermissionLevel(ctx, guild.Id, userId)
			if err != nil {
				return err
			}

			mu.Lock()
			dtos = append(dtos, GuildDto{
				Id:              guild.Id,
				Name:            guild.Name,
				Icon:            guild.Icon,
				PermissionLevel: permLevel,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	// Sort the guilds by name, but put the guilds with permission_level=0 last
	slices.SortFunc(dtos, func(a, b GuildDto) int {
		if a.PermissionLevel == 0 && b.PermissionLevel > 0 {
			return 1
		} else if a.PermissionLevel > 0 && b.PermissionLevel == 0 {
			return -1
		}

		return cmp.Compare(a.Name, b.Name)
	})

	return dtos, nil
}

// TODO: Remove this function!
func storeGuildsInDb(ctx context.Context, userId uint64, guilds []guild.Guild) error {
	var wrappedGuilds []database.UserGuild

	// endpoint's partial guild doesn't includes ownerid
	// we only user cached guilds on the index page, so it doesn't matter if we don't have have the real owner id
	// if the user isn't the owner, as we pull from the cache on other endpoints
	for _, guild := range guilds {
		wrappedGuilds = append(wrappedGuilds, database.UserGuild{
			GuildId:         guild.Id,
			Name:            guild.Name,
			Owner:           guild.Owner,
			UserPermissions: guild.Permissions,
			Icon:            guild.Icon,
		})
	}

	return dbclient.Client.UserGuilds.Set(ctx, userId, wrappedGuilds)
}

func getGuildIntersection(ctx context.Context, userId uint64, userGuilds []guild.Guild) ([]guild.Guild, error) {
	guildIds := make([]uint64, len(userGuilds))
	for i, guild := range userGuilds {
		guildIds[i] = guild.Id
	}

	// Restrict the set of guilds to guilds that the bot is also in
	botGuilds, err := getExistingGuilds(ctx, guildIds)
	if err != nil {
		return nil, err
	}

	botGuildIds := collections.NewSet[uint64]()
	for _, guildId := range botGuilds {
		botGuildIds.Add(guildId)
	}

	// Get the intersection of the two sets
	intersection := make([]guild.Guild, 0, len(botGuilds))
	for _, guild := range userGuilds {
		if botGuildIds.Contains(guild.Id) {
			intersection = append(intersection, guild)
		}
	}

	return intersection, nil
}

func getExistingGuilds(ctx context.Context, userGuilds []uint64) ([]uint64, error) {
	query := `SELECT "guild_id" from guilds WHERE "guild_id" = ANY($1);`

	userGuildsArray := &pgtype.Int8Array{}
	if err := userGuildsArray.Set(userGuilds); err != nil {
		return nil, err
	}

	rows, err := cache.Instance.Query(ctx, query, userGuildsArray)
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

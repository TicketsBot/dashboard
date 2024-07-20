package api

import (
	"cmp"
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/collections"
	"github.com/TicketsBot/common/permission"
	syncutils "github.com/TicketsBot/common/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"golang.org/x/sync/errgroup"
	"slices"
	"sync"
	"time"
)

type wrappedGuild struct {
	Id              uint64                     `json:"id,string"`
	Name            string                     `json:"name"`
	Icon            string                     `json:"icon"`
	PermissionLevel permission.PermissionLevel `json:"permission_level"`
}

func GetGuilds(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	// Get the guilds that the user is in, that the bot is also in
	userGuilds, err := getGuildIntersection(userId)
	if err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	wg := syncutils.NewChannelWaitGroup()
	wg.Add(len(userGuilds))

	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex
	guilds := make([]wrappedGuild, 0, len(userGuilds))
	for _, guild := range userGuilds {
		guild := guild

		group.Go(func() error {
			defer wg.Done()

			permLevel, err := utils.GetPermissionLevel(ctx, guild.GuildId, userId)
			if err != nil {
				return err
			}

			mu.Lock()
			guilds = append(guilds, wrappedGuild{
				Id:              guild.GuildId,
				Name:            guild.Name,
				Icon:            guild.Icon,
				PermissionLevel: permLevel,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// Sort the guilds by name, but put the guilds with permission_level=0 last
	slices.SortFunc(guilds, func(a, b wrappedGuild) int {
		if a.PermissionLevel == 0 && b.PermissionLevel > 0 {
			return 1
		} else if a.PermissionLevel > 0 && b.PermissionLevel == 0 {
			return -1
		}

		return cmp.Compare(a.Name, b.Name)
	})

	c.JSON(200, guilds)
}

func getGuildIntersection(userId uint64) ([]database.UserGuild, error) {
	// Get all the guilds that the user is in
	userGuilds, err := dbclient.Client.UserGuilds.Get(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	guildIds := make([]uint64, len(userGuilds))
	for i, guild := range userGuilds {
		guildIds[i] = guild.GuildId
	}

	// Restrict the set of guilds to guilds that the bot is also in
	botGuilds, err := getExistingGuilds(guildIds)
	if err != nil {
		return nil, err
	}

	botGuildIds := collections.NewSet[uint64]()
	for _, guildId := range botGuilds {
		botGuildIds.Add(guildId)
	}

	// Get the intersection of the two sets
	intersection := make([]database.UserGuild, 0, len(botGuilds))
	for _, guild := range userGuilds {
		if botGuildIds.Contains(guild.GuildId) {
			intersection = append(intersection, guild)
		}
	}

	return intersection, nil
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

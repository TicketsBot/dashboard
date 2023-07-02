package cache

import (
	"context"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/jackc/pgx/v4/pgxpool"
	gdlcache "github.com/rxdn/gdl/cache"
)

type Cache struct {
	*gdlcache.PgCache
}

var Instance *Cache

func NewCache() *Cache {
	pool, err := pgxpool.Connect(context.Background(), config.Conf.Cache.Uri)
	if err != nil {
		panic(err)
	}

	cache := gdlcache.NewPgCache(pool, gdlcache.CacheOptions{
		Guilds:   true,
		Users:    true,
		Members:  true,
		Channels: true,
		Roles:    false,
	})

	return &Cache{
		PgCache: &cache,
	}
}

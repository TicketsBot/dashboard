package redis

import (
	"context"
	"fmt"
	"time"
)

const ChannelRefreshCooldown = 60 * time.Second

func (c *RedisClient) TakeChannelRefreshToken(ctx context.Context, guildId uint64) (bool, error) {
	key := fmt.Sprintf("tickets:channelrefershcooldown:%d", guildId)

	res, err := c.SetNX(ctx, key, "1", ChannelRefreshCooldown).Result()
	if err != nil {
		return false, err
	}

	return res, nil
}

package botcontext

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	cacheclient "github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/restcache"
	"github.com/TicketsBot/database"
	cache "github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type BotContext struct {
	BotId       uint64
	Token       string
	RateLimiter *ratelimit.Ratelimiter
	RestCache   restcache.RestCache
}

func (c *BotContext) Db() *database.Database {
	return dbclient.Client
}

func (c *BotContext) Cache() permission.PermissionCache {
	return permission.NewRedisCache(redis.Client.Client)
}

func (c *BotContext) IsBotAdmin(_ context.Context, userId uint64) bool {
	for _, id := range config.Conf.Admins {
		if id == userId {
			return true
		}
	}

	return false
}

func (c *BotContext) GetGuild(ctx context.Context, guildId uint64) (guild.Guild, error) {
	g, err := cacheclient.Instance.GetGuild(ctx, guildId)
	switch {
	case err == nil:
		return g, nil
	case errors.Is(err, cache.ErrNotFound):
		g, err := rest.GetGuild(ctx, c.Token, c.RateLimiter, guildId)
		if err == nil {
			if err := cacheclient.Instance.StoreGuild(ctx, g); err != nil {
				return guild.Guild{}, err
			}
		}

		return g, err
	default:
		return guild.Guild{}, err
	}
}

func (c *BotContext) GetGuildOwner(ctx context.Context, guildId uint64) (uint64, error) {
	cachedOwner, err := cacheclient.Instance.GetGuildOwner(ctx, guildId)
	switch {
	case err == nil:
		return cachedOwner, nil
	case errors.Is(err, cache.ErrNotFound):
		guild, err := c.GetGuild(ctx, guildId)
		if err != nil {
			return 0, err
		}

		if err := cacheclient.Instance.StoreGuild(ctx, guild); err != nil {
			return 0, err
		}

		return guild.OwnerId, nil
	default:
		return 0, err
	}
}

func (c *BotContext) GetGuildMember(ctx context.Context, guildId, userId uint64) (member.Member, error) {
	m, err := cacheclient.Instance.GetMember(ctx, guildId, userId)
	switch {
	case err == nil:
		return m, nil
	case errors.Is(err, cache.ErrNotFound):
		m, err := rest.GetGuildMember(ctx, c.Token, c.RateLimiter, guildId, userId)
		if err != nil {
			return member.Member{}, nil
		}

		if err := cacheclient.Instance.StoreMember(ctx, m, guildId); err != nil {
			return member.Member{}, err
		}

		return m, nil
	default:
		return member.Member{}, err
	}
}

func (c *BotContext) RemoveGuildMemberRole(ctx context.Context, guildId, userId, roleId uint64) error {
	return rest.RemoveGuildMemberRole(ctx, c.Token, c.RateLimiter, guildId, userId, roleId)
}

func (c *BotContext) CreateGuildRole(ctx context.Context, guildId uint64, data rest.GuildRoleData) (guild.Role, error) {
	return rest.CreateGuildRole(ctx, c.Token, c.RateLimiter, guildId, data)
}

func (c *BotContext) DeleteGuildRole(ctx context.Context, guildId, roleId uint64) error {
	return rest.DeleteGuildRole(ctx, c.Token, c.RateLimiter, guildId, roleId)
}

func (c *BotContext) GetUser(ctx context.Context, userId uint64) (user.User, error) {
	u, err := cacheclient.Instance.GetUser(ctx, userId)
	switch {
	case err == nil:
		return u, nil
	case errors.Is(err, cache.ErrNotFound):
		u, err := rest.GetUser(ctx, c.Token, c.RateLimiter, userId)
		if err != nil {
			return user.User{}, err
		}

		if err := cacheclient.Instance.StoreUser(ctx, u); err != nil {
			return user.User{}, err
		}

		return u, nil
	default:
		return user.User{}, err
	}
}

func (c *BotContext) GetGuildRoles(ctx context.Context, guildId uint64) ([]guild.Role, error) {
	return c.RestCache.GetGuildRoles(guildId)
}

func (c *BotContext) GetGuildChannels(ctx context.Context, guildId uint64) ([]channel.Channel, error) {
	cachedChannels, err := cacheclient.Instance.GetGuildChannels(ctx, guildId)
	if err != nil {
		return nil, err
	}

	if len(cachedChannels) > 0 {
		return cachedChannels, nil
	} else {
		// If guild is cached but not any channels, likely that it does truly have 0 channels,
		// so don't fetch from REST.
		_, err := cacheclient.Instance.GetGuild(ctx, guildId)
		switch {
		case err == nil:
			return []channel.Channel{}, nil
		case errors.Is(err, cache.ErrNotFound):
			// If guild isn't cached, fetch from REST
			channels, err := rest.GetGuildChannels(ctx, c.Token, c.RateLimiter, guildId)
			if err != nil {
				return nil, err
			}

			if err := cacheclient.Instance.StoreChannels(ctx, channels); err != nil {
				return nil, err
			}

			return channels, nil
		default:
			return nil, err
		}
	}
}

func (c *BotContext) GetGuildEmoji(ctx context.Context, guildId, emojiId uint64) (emoji.Emoji, error) {
	e, err := cacheclient.Instance.GetEmoji(ctx, guildId)
	switch {
	case err == nil:
		return e, nil
	case errors.Is(err, cache.ErrNotFound):
		e, err := rest.GetGuildEmoji(ctx, c.Token, c.RateLimiter, guildId, emojiId)
		if err != nil {
			return emoji.Emoji{}, err
		}

		if err := cacheclient.Instance.StoreEmoji(ctx, e, guildId); err != nil {
			return emoji.Emoji{}, err
		}

		return e, nil
	default:
		return emoji.Emoji{}, err
	}
}

func (c *BotContext) GetGuildEmojis(ctx context.Context, guildId uint64) ([]emoji.Emoji, error) {
	emojis, err := cacheclient.Instance.GetGuildEmojis(ctx, guildId)
	if err != nil {
		return nil, err
	}

	if len(emojis) == 0 {
		emojis, err := rest.ListGuildEmojis(ctx, c.Token, c.RateLimiter, guildId)
		if err != nil {
			return nil, err
		}

		if err := cacheclient.Instance.StoreEmojis(ctx, emojis, guildId); err != nil {
			return nil, err
		}

		return emojis, err
	}

	return emojis, nil
}

func (c *BotContext) SearchMembers(ctx context.Context, guildId uint64, query string) ([]member.Member, error) {
	data := rest.SearchGuildMembersData{
		Query: query,
		Limit: 100,
	}

	members, err := rest.SearchGuildMembers(ctx, c.Token, c.RateLimiter, guildId, data)
	if err != nil {
		return nil, err
	}

	if err := cacheclient.Instance.StoreMembers(ctx, members, guildId); err != nil {
		return nil, err
	}

	return members, nil
}

func (c *BotContext) ListMembers(ctx context.Context, guildId uint64) ([]member.Member, error) {
	data := rest.ListGuildMembersData{
		Limit: 100,
	}

	members, err := rest.ListGuildMembers(ctx, c.Token, c.RateLimiter, guildId, data)
	if err != nil {
		return nil, err
	}

	if err := cacheclient.Instance.StoreMembers(ctx, members, guildId); err != nil {
		return nil, err
	}

	return members, nil
}

func (c *BotContext) CreateGuildCommand(ctx context.Context, guildId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGuildCommand(ctx, c.Token, c.RateLimiter, c.BotId, guildId, data)
}

func (c *BotContext) DeleteGuildCommand(ctx context.Context, guildId, commandId uint64) error {
	return rest.DeleteGuildCommand(ctx, c.Token, c.RateLimiter, c.BotId, guildId, commandId)
}

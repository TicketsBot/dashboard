package botcontext

import (
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/restcache"
	"github.com/TicketsBot/database"
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

func (ctx BotContext) Db() *database.Database {
	return dbclient.Client
}

func (ctx BotContext) Cache() permission.PermissionCache {
	return permission.NewRedisCache(redis.Client.Client)
}

func (ctx BotContext) IsBotAdmin(userId uint64) bool {
	for _, id := range config.Conf.Admins {
		if id == userId {
			return true
		}
	}

	return false
}

func (ctx BotContext) GetGuild(guildId uint64) (g guild.Guild, err error) {
	if guild, found := cache.Instance.GetGuild(guildId); found {
		return guild, nil
	}

	g, err = rest.GetGuild(ctx.Token, ctx.RateLimiter, guildId)
	if err == nil {
		go cache.Instance.StoreGuild(g)
	}

	return
}

func (ctx BotContext) GetGuildOwner(guildId uint64) (uint64, error) {
	cachedOwner, exists := cache.Instance.GetGuildOwner(guildId)
	if exists {
		return cachedOwner, nil
	}

	guild, err := ctx.GetGuild(guildId)
	if err != nil {
		return 0, err
	}

	go cache.Instance.StoreGuild(guild)
	return guild.OwnerId, nil
}

func (ctx BotContext) GetGuildMember(guildId, userId uint64) (m member.Member, err error) {
	if guild, found := cache.Instance.GetMember(guildId, userId); found {
		return guild, nil
	}

	m, err = rest.GetGuildMember(ctx.Token, ctx.RateLimiter, guildId, userId)
	if err == nil {
		go cache.Instance.StoreMember(m, guildId)
	}

	return
}

func (ctx BotContext) RemoveGuildMemberRole(guildId, userId, roleId uint64) (err error) {
	err = rest.RemoveGuildMemberRole(ctx.Token, ctx.RateLimiter, guildId, userId, roleId)
	return
}

func (ctx BotContext) CreateGuildRole(guildId uint64, data rest.GuildRoleData) (role guild.Role, err error) {
	role, err = rest.CreateGuildRole(ctx.Token, ctx.RateLimiter, guildId, data)
	return
}

func (ctx BotContext) DeleteGuildRole(guildId, roleId uint64) (err error) {
	err = rest.DeleteGuildRole(ctx.Token, ctx.RateLimiter, guildId, roleId)
	return
}

func (ctx BotContext) GetUser(userId uint64) (u user.User, err error) {
	u, err = rest.GetUser(ctx.Token, ctx.RateLimiter, userId)
	if err == nil {
		go cache.Instance.StoreUser(u)
	}

	return
}

func (ctx BotContext) GetGuildRoles(guildId uint64) (roles []guild.Role, err error) {
	return ctx.RestCache.GetGuildRoles(guildId)
}

func (ctx BotContext) GetGuildEmoji(guildId, emojiId uint64) (emoji.Emoji, error) {
	if emoji, ok := cache.Instance.GetEmoji(guildId); ok {
		return emoji, nil
	}

	emoji, err := rest.GetGuildEmoji(ctx.Token, ctx.RateLimiter, guildId, emojiId)
	if err == nil {
		go cache.Instance.StoreEmoji(emoji, guildId)
	}

	return emoji, err
}

func (ctx BotContext) GetGuildEmojis(guildId uint64) (emojis []emoji.Emoji, err error) {
	if emojis := cache.Instance.GetGuildEmojis(guildId); len(emojis) > 0 {
		return emojis, nil
	}

	emojis, err = rest.ListGuildEmojis(ctx.Token, ctx.RateLimiter, guildId)
	if err == nil {
		go cache.Instance.StoreEmojis(emojis, guildId)
	}

	return
}

func (ctx BotContext) SearchMembers(guildId uint64, query string) (members []member.Member, err error) {
	data := rest.SearchGuildMembersData{
		Query: query,
		Limit: 100,
	}

	members, err = rest.SearchGuildMembers(ctx.Token, ctx.RateLimiter, guildId, data)
	if err == nil {
		go cache.Instance.StoreMembers(members, guildId)
	}

	return
}

func (ctx BotContext) ListMembers(guildId uint64) (members []member.Member, err error) {
	data := rest.ListGuildMembersData{
		Limit: 100,
	}

	members, err = rest.ListGuildMembers(ctx.Token, ctx.RateLimiter, guildId, data)
	if err == nil {
		go cache.Instance.StoreMembers(members, guildId)
	}

	return
}

func (ctx BotContext) CreateGuildCommand(guildId uint64, data rest.CreateCommandData) (interaction.ApplicationCommand, error) {
	return rest.CreateGuildCommand(ctx.Token, ctx.RateLimiter, ctx.BotId, guildId, data)
}

func (ctx BotContext) DeleteGuildCommand(guildId, commandId uint64) error {
	return rest.DeleteGuildCommand(ctx.Token, ctx.RateLimiter, ctx.BotId, guildId, commandId)
}

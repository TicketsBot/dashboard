package botcontext

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func ContextForGuild(guildId uint64) (ctx BotContext, err error) {
	whitelabelBotId, isWhitelabel, err := dbclient.Client.WhitelabelGuilds.GetBotByGuild(guildId)
	if err != nil {
		return
	}

	var keyPrefix string

	if isWhitelabel {
		res, err := dbclient.Client.Whitelabel.GetByBotId(whitelabelBotId)
		if err != nil {
			return ctx, err
		}

		ctx.BotId = res.BotId
		ctx.Token = res.Token
		keyPrefix = fmt.Sprintf("ratelimiter:%d", whitelabelBotId)
	} else {
		ctx.BotId = config.Conf.Bot.Id
		ctx.Token = config.Conf.Bot.Token
		keyPrefix = "ratelimiter:public"
	}

	// TODO: Large sharding buckets
	ctx.RateLimiter = ratelimit.NewRateLimiter(ratelimit.NewRedisStore(redis.Client.Client, keyPrefix), 1)

	return
}

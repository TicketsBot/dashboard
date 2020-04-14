package ratelimit

import (
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/rxdn/gdl/rest/ratelimit"
)

var Ratelimiter = ratelimit.NewRateLimiter(ratelimit.NewRedisStore(messagequeue.Client.Client, "ratelimit")) // TODO: Use values from config

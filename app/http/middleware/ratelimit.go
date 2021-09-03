package middleware

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"hash/fnv"
	"strconv"
	"time"
)

type RateLimitType uint8

const (
	RateLimitTypeIp RateLimitType = iota
	RateLimitTypeUser
	RateLimitTypeGuild
)

func CreateRateLimiter(rlType RateLimitType, max int, period time.Duration) gin.HandlerFunc {
	limiter := redis_rate.NewLimiter(redis.Client)

	return func(ctx *gin.Context) {
		limit := redis_rate.Limit{
			Rate:   max,
			Burst:  max,
			Period: period,
		}

		name, skip := getKey(ctx, rlType, limit)
		if skip {
			ctx.Next()
			return
		}

		res, err := limiter.Allow(redis.DefaultContext(), name, limit)
		if err != nil {
			ctx.AbortWithStatusJSON(500, utils.ErrorJson(err))
			Logging(ctx)
			return
		}

		// Use smallest remaining for ratelimit headers
		smallestRemaining := ctx.Keys["rl_sr"]
		if smallestRemaining == nil {
			writeHeaders(ctx, res)
		} else {
			rem := smallestRemaining.(int)
			if res.Remaining < rem {
				writeHeaders(ctx, res)
			}
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(res.Limit.Rate))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Header("X-RateLimit-Reset-After", strconv.FormatInt(res.ResetAfter.Milliseconds(), 10))

		if res.Allowed <= 0 {
			ctx.AbortWithStatusJSON(429, utils.ErrorStr("You are being ratelimited"))
			Logging(ctx)
			return
		}

		ctx.Next()
	}
}

func writeHeaders(ctx *gin.Context, res *redis_rate.Result) {
	ctx.Keys["rl_sr"] = res.Remaining
	fmt.Println(res.Remaining)
	ctx.Header("X-RateLimit-Limit", strconv.Itoa(res.Limit.Rate))
	ctx.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
	ctx.Header("X-RateLimit-Reset-After", strconv.FormatInt(res.ResetAfter.Milliseconds(), 10))
}

// Returns (key, skip)
func getKey(ctx *gin.Context, rlType RateLimitType, limit redis_rate.Limit) (string, bool) {
	userId := ctx.Keys["userid"]
	guildId := ctx.Keys["guildid"]

	if (rlType == RateLimitTypeUser && userId == nil) || (rlType == RateLimitTypeGuild && guildId == nil) {
		ctx.Next()
		return "", true
	}

	var key string
	switch rlType {
	case RateLimitTypeIp:
		key = ctx.ClientIP()
	case RateLimitTypeUser:
		key = strconv.FormatUint(userId.(uint64), 10)
	case RateLimitTypeGuild:
		key = strconv.FormatUint(guildId.(uint64), 10)
	}

	target := fmt.Sprintf("%d:%s", rlType, key)
	bucket := fmt.Sprintf("%d/%d", limit.Rate, limit.Period.Milliseconds())
	full := fmt.Sprintf("%s:%s:%s", target, bucket, ctx.FullPath())

	return strconv.FormatUint(uint64(hash(full)), 16), false
}

func hash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

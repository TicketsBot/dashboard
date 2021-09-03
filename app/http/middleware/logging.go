package middleware

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
)

func Logging(ctx *gin.Context) {
	ctx.Next()

	statusCode := ctx.Writer.Status()

	if !config.Conf.Debug && statusCode >= 200 && statusCode <= 299 {
		return
	}

	level := sentry.LevelInfo
	if statusCode >= 500 {
		level = sentry.LevelError
	}

	body, _ := ioutil.ReadAll(ctx.Request.Body)

	sentry.CaptureEvent(&sentry.Event{
		Extra: map[string]interface{}{
			"status_code": strconv.Itoa(statusCode),
			"method": ctx.Request.Method,
			"path": ctx.Request.URL.Path,
			"guild_id": ctx.Keys["guildid"],
			"user_id": ctx.Keys["userid"],
			"body": string(body),
		},
		Level:      level,
		Message:    fmt.Sprintf("HTTP %d on %s %s", statusCode, ctx.Request.Method, ctx.FullPath()),
		Tags: map[string]string{
			"status_code": strconv.Itoa(statusCode),
			"method": ctx.Request.Method,
			"path": ctx.Request.URL.Path,
		},
	})
}

package middleware

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
)

func Logging(ctx *gin.Context) {
	ctx.Next()

	statusCode := ctx.Writer.Status()

	var level sentry.Level
	if statusCode >= 500 {
		level = sentry.LevelError
	} else if statusCode >= 400 {
		level = sentry.LevelWarning
	} else {
		level = sentry.LevelInfo
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

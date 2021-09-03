package middleware

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
)

func Logging(ctx *gin.Context) {
	defer ctx.Next()

	statusCode := ctx.Writer.Status()

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
			"user_id": ctx.Keys["user_id"],
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

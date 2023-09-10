package middleware

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"runtime/debug"
	"strconv"
)

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

func (l Level) sentryLevel() sentry.Level {
	switch l {
	case LevelDebug:
		return sentry.LevelDebug
	case LevelInfo:
		return sentry.LevelInfo
	case LevelWarning:
		return sentry.LevelWarning
	case LevelError:
		return sentry.LevelError
	case LevelFatal:
		return sentry.LevelFatal
	default:
		return sentry.LevelDebug
	}
}

func Logging(minLevel Level) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		statusCode := ctx.Writer.Status()

		level := LevelInfo
		if statusCode >= 500 {
			level = LevelError
		} else if statusCode >= 400 {
			level = LevelWarning
		}

		if level < minLevel {
			return
		}

		requestBody, _ := ioutil.ReadAll(ctx.Request.Body)

		var responseBody []byte
		if statusCode >= 400 && statusCode <= 599 {
			cw, ok := ctx.Writer.(*CustomWriter)
			if ok {
				responseBody = cw.Read()
			}
		}

		sentry.CaptureEvent(&sentry.Event{
			Extra: map[string]interface{}{
				"status_code":  strconv.Itoa(statusCode),
				"method":       ctx.Request.Method,
				"path":         ctx.Request.URL.Path,
				"query":        ctx.Request.URL.RawQuery,
				"guild_id":     ctx.Keys["guildid"],
				"user_id":      ctx.Keys["userid"],
				"request_body": string(requestBody),
				"response":     string(responseBody),
				"stacktrace":   string(debug.Stack()),
			},
			Level:   level.sentryLevel(),
			Message: fmt.Sprintf("HTTP %d on %s %s", statusCode, ctx.Request.Method, ctx.FullPath()),
			Tags: map[string]string{
				"status_code": strconv.Itoa(statusCode),
				"method":      ctx.Request.Method,
				"path":        ctx.Request.URL.Path,
			},
		})
	}
}

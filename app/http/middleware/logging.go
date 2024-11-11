package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		statusCode := c.Writer.Status()

		level := zapcore.InfoLevel
		if statusCode >= 500 {
			level = zapcore.ErrorLevel
		} else if statusCode >= 400 {
			level = zapcore.WarnLevel
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", c.Writer.Status()),
			zap.String("timestamp", start.String()),
			zap.Duration("latency", time.Now().Sub(start)),
			zap.String("client_ip", c.ClientIP()),
		}

		if guildId, ok := c.Keys["guildid"]; ok {
			fields = append(fields, zap.Uint64("guild_id", guildId.(uint64)))
		}

		if userId, ok := c.Keys["userid"]; ok {
			fields = append(fields, zap.Uint64("user_id", userId.(uint64)))
		}

		logger.Log(level, "Incoming HTTP request", fields...)
	}
}

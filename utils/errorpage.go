package utils

import "github.com/gin-gonic/gin"

func ErrorPage(ctx *gin.Context, statusCode int, error string) {
	ctx.HTML(statusCode, "error", gin.H{
		"status": statusCode,
		"error": error,
	})
}

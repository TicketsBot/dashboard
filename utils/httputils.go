package utils

import "github.com/gin-gonic/gin"

func Respond(ctx *gin.Context, s string) {
	ctx.Data(200, "text/html; charset=utf-8", []byte(s))
}

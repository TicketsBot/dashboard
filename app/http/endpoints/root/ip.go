package root

import "github.com/gin-gonic/gin"

func IpHandler(ctx *gin.Context)  {
	ctx.String(200, ctx.ClientIP())
}

package middleware

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseGuildId(ctx *gin.Context) {
	guildId, ok := ctx.Params.Get("id")
	if !ok {
		ctx.AbortWithStatusJSON(400, utils.ErrorStr("Missing guild ID"))
		return
	}

	parsed, err := strconv.ParseUint(guildId, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, utils.ErrorStr("Invalid guild ID"))
		return
	}

	ctx.Keys["guildid"] = parsed
}

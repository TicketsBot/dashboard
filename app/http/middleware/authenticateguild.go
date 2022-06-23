package middleware

import (
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"strconv"
)

// requires AuthenticateCookie middleware to be run before
func AuthenticateGuild(requiredPermissionLevel permission.PermissionLevel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if guildId, ok := ctx.Params.Get("id"); ok {
			parsed, err := strconv.ParseUint(guildId, 10, 64)
			if err != nil {
				ctx.JSON(400, utils.ErrorStr("Invalid guild ID"))
				ctx.Abort()
				return
			}

			ctx.Keys["guildid"] = parsed

			// TODO: Do we need this? Only really serves as a check whether the bot is in the server
			guild, found := cache.Instance.GetGuild(parsed, false)
			if !found {
				ctx.JSON(404, utils.ErrorStr("Guild not found"))
				ctx.Abort()
				return
			}

			// Verify the user has permissions to be here
			userId := ctx.Keys["userid"].(uint64)

			permLevel, err := utils.GetPermissionLevel(guild.Id, userId)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				ctx.Abort()
				return
			}

			if permLevel < requiredPermissionLevel {
				ctx.JSON(403, utils.ErrorStr("Unauthorized"))
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(400, utils.ErrorStr("Invalid guild ID"))
			ctx.Abort()
			return
		}
	}
}

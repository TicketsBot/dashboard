package middleware

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"strconv"
)

// requires AuthenticateCookie middleware to be run before
func AuthenticateGuild(isApiMethod bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if guildId, ok := ctx.Params.Get("id"); ok {
			parsed, err := strconv.ParseUint(guildId, 10, 64)
			if err != nil {
				if isApiMethod {
					ctx.JSON(400, gin.H{
						"success": false,
						"error": "Invalid guild ID",
					})
				} else {
					utils.ErrorPage(ctx, 400, "Invalid server ID")
					ctx.Abort()
				}
				ctx.Abort()
				return
			}

			ctx.Keys["guildid"] = parsed

			guild, found := cache.Instance.GetGuild(parsed, false)
			if !found {
				if isApiMethod {
					ctx.Redirect(302, fmt.Sprintf("https://invite.ticketsbot.net/?guild_id=%d&disable_guild_select=true&response_type=code&scope=bot%%20identify&redirect_uri=%s", parsed, config.Conf.Server.BaseUrl))
				} else {
					utils.ErrorPage(ctx, 404, "Couldn't find server with the provided ID")
				}
				ctx.Abort()
				return
			}

			ctx.Keys["guild"] = guild

			// Verify the user has permissions to be here
			userId := ctx.Keys["userid"].(uint64)
			if utils.GetPermissionLevel(guild.Id, userId) != permission.Admin {
				if isApiMethod {
					ctx.AbortWithStatusJSON(403, gin.H{
						"success": false,
						"error": "Unauthorized",
					})
				} else {
					utils.ErrorPage(ctx, 403, "You don't have permission to be here - make sure the server owner has added you as an admin with t!addadmin")
					ctx.Abort()
				}
			}
		} else {
			if isApiMethod {
				ctx.JSON(400, gin.H{
					"success": false,
					"error": "Invalid guild ID",
				})
			} else {
				utils.ErrorPage(ctx, 400, "Invalid server ID")
			}
			ctx.Abort()
		}
	}
}

package root

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userIdStr := store.Get("userid").(string)
		userId, err := utils.GetUserId(store)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		userGuilds := table.GetGuilds(userIdStr)
		adminGuilds := make([]objects.Guild, 0)
		for _, guild := range userGuilds {
			guildId, err := strconv.ParseInt(guild.Id, 10, 64)
			if err != nil { // I think this happens when a server was deleted? We should just skip though.
				continue
			}

			isAdmin := make(chan bool)
			go utils.IsAdmin(guild, guildId, userId, isAdmin)
			if <-isAdmin {
				adminGuilds = append(adminGuilds, guild)
			}
		}

		ctx.HTML(200, "main/index", gin.H{
			"name":    store.Get("name").(string),
			"baseurl": config.Conf.Server.BaseUrl,
			"servers": adminGuilds,
			"empty":   len(adminGuilds) == 0,
			"isIndex": true,
			"avatar": store.Get("avatar").(string),
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

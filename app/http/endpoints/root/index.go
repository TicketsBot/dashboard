package root

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
	"strconv"
)

func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userId, err := utils.GetUserId(store)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		userGuilds := table.GetGuilds(userId)
		adminGuilds := make([]objects.Guild, 0)
		for _, g := range userGuilds {
			guildId, err := strconv.ParseUint(g.Id, 10, 64)
			if err != nil { // I think this happens when a server was deleted? We should just skip though.
				continue
			}

			fakeGuild := guild.Guild{
				Owner:       g.Owner,
				Permissions: g.Permissions,
			}

			isAdmin := make(chan bool)
			go utils.IsAdmin(fakeGuild, guildId, userId, isAdmin)
			if <-isAdmin {
				adminGuilds = append(adminGuilds, g)
			}
		}

		ctx.HTML(200, "main/index", gin.H{
			"name":    store.Get("name").(string),
			"baseurl": config.Conf.Server.BaseUrl,
			"servers": adminGuilds,
			"empty":   len(adminGuilds) == 0,
			"isIndex": true,
			"avatar":  store.Get("avatar").(string),
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

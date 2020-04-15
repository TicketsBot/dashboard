package root

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
)
0
func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userId := utils.GetUserId(store)

		userGuilds := table.GetGuilds(userId)
		adminGuilds := make([]guild.Guild, 0)
		for _, g := range userGuilds {
			fakeGuild := guild.Guild{
				Id:          g.Id,
				OwnerId:     g.OwnerId,
				Permissions: g.Permissions,
			}

			isAdmin := make(chan bool)
			go utils.IsAdmin(fakeGuild, userId, isAdmin)
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

package root

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
)

func IndexHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	userId := utils.GetUserId(store)

	userGuilds, err := database.Client.UserGuilds.Get(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	adminGuilds := make([]guild.Guild, 0)
	for _, g := range userGuilds {
		fakeGuild := guild.Guild{
			Id:          g.GuildId,
			Owner:       g.Owner,
			Permissions: int(g.UserPermissions),
		}

		if g.Owner {
			fakeGuild.OwnerId = userId
		}

		isAdmin := make(chan bool)
		go utils.IsAdmin(fakeGuild, userId, isAdmin)
		if <-isAdmin {
			adminGuilds = append(adminGuilds, fakeGuild)
		}
	}

	ctx.HTML(200, "main/index", gin.H{
		"name":    store.Get("name").(string),
		"baseurl": config.Conf.Server.BaseUrl,
		"avatar":  store.Get("avatar").(string),
	})
}

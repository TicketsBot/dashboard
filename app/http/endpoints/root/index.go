package root

import (
	"github.com/TicketsBot/GoPanel/app/http/template"
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
		userId, err := utils.GetUserId(store); if err != nil {
			ctx.String(500, err.Error())
			return
		}

		adminGuilds := make([]objects.Guild, 0)
		adminGuildIds := table.GetAdminGuilds(userId)
		for _, guild := range table.GetGuilds(userIdStr) {
			guildId, err := strconv.ParseInt(guild.Id, 10, 64); if err != nil {
				ctx.String(500, err.Error())
				return
			}

			if guild.Owner || utils.Contains(adminGuildIds, guildId) {
				adminGuilds = append(adminGuilds, guild)
			}
		}

		var servers []map[string]string
		for _, server := range adminGuilds {
			element := map[string]string{
				"serverid": server.Id,
				"servername": server.Name,
			}

			servers = append(servers, element)
		}

		utils.Respond(ctx, template.TemplateIndex.Render(map[string]interface{}{
			"name": store.Get("name").(string),
			"baseurl": config.Conf.Server.BaseUrl,
			"servers": servers,
			"empty": len(servers) == 0,
		}))
	} else {
		ctx.Redirect(302, "/login")
	}
}


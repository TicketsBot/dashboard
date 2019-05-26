package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/template"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func BlacklistHandler(ctx *gin.Context) {
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

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseInt(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		var guild objects.Guild
		for _, g := range table.GetGuilds(userIdStr) {
			if g.Id == guildIdStr {
				guild = g
				break
			}
		}

		// Verify the user has permissions to be here
		if !guild.Owner && !table.IsAdmin(guildId, userId) {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		curid, err := strconv.ParseInt("210105426849562625", 10, 64)
		table.RemoveBlacklist(guildId, curid)
		fmt.Println(table.IsBlacklisted(guildId, curid))

		//blacklistedUsers := table.GetBlacklistNodes(guildId)
		//var blacklisted []map[string]interface{}
		//for _, node := range blacklistedUsers {

		//}

		fmt.Println(len(guild.Members))

		utils.Respond(ctx, template.TemplateBlacklist.Render(map[string]interface{}{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
		}))
	}
}

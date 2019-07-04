package manage

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

		blacklistedUsers := table.GetBlacklistNodes(guildId)

		var blacklistedIds []int64
		for _, user := range blacklistedUsers {
			blacklistedIds = append(blacklistedIds, user.User)
		}

		nodes := table.GetUserNodes(blacklistedIds)

		var blacklisted []map[string]interface{}
		for _, node := range nodes {
			blacklisted = append(blacklisted, map[string]interface{}{
				"userId": node.Id,
				"username": utils.Base64Decode(node.Name),
				"discrim": node.Discriminator,
			})
		}

		userNotFound := false
		isStaff := false
		if store.Get("csrf").(string) == ctx.Query("csrf") { // CSRF is correct *and* set
			username := ctx.Query("username")
			discrim := ctx.Query("discrim")

			// Verify that the user ID is real and in a shared guild
			targetId := table.GetUserId(username, discrim)
			exists := targetId != 0

			if exists {
				if guild.OwnerId == strconv.Itoa(int(targetId)) || table.IsStaff(guildId, targetId) { // Prevent users from blacklisting staff
					isStaff = true
				} else {
					if !utils.Contains(blacklistedIds, targetId) { // Prevent duplicates
						table.AddBlacklist(guildId, targetId)
						blacklisted = append(blacklisted, map[string]interface{}{
							"userId": targetId,
							"username": username,
							"discrim": discrim,
						})
					}
				}
			} else {
				userNotFound = true
			}
		}

		utils.Respond(ctx, template.TemplateBlacklist.Render(map[string]interface{}{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"csrf": store.Get("csrf").(string),
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
			"blacklisted": blacklisted,
			"userNotFound": userNotFound,
			"isStaff": isStaff,
		}))
	}
}

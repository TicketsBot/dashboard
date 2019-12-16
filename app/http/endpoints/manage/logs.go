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

func LogsHandler(ctx *gin.Context) {
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

		pageStr := ctx.Param("page")
		page := 1
		i, err := strconv.Atoi(pageStr)
		if err == nil {
			if i > 0 {
				page = i
			}
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
		if !utils.Contains(config.Conf.Admins, userIdStr) && !guild.Owner && !table.IsAdmin(guildId, userId) {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		pageLimit := 30

		// Get logs
		// Get user ID from URL
		var filteredUserId int64
		if utils.IsInt(ctx.Query("userid")) {
			filteredUserId, _ = strconv.ParseInt(ctx.Query("userid"), 10, 64)
		}

		// Get ticket ID from URL
		var ticketId int
		if utils.IsInt(ctx.Query("ticketid")) {
			ticketId, _ = strconv.Atoi(ctx.Query("ticketid"))
		}

		// Get username from URL
		username := ctx.Query("username")

		// Get logs from DB
		logs := table.GetFilteredTicketArchives(guildId, filteredUserId, username, ticketId)

		// Select 30 logs + format them
		var formattedLogs []map[string]interface{}
		for i := (page - 1) * pageLimit; i < (page - 1) * pageLimit + pageLimit; i++ {
			if i >= len(logs) {
				break
			}

			log := logs[i]
			formattedLogs = append(formattedLogs, map[string]interface{}{
				"ticketid":  log.TicketId,
				"userid": log.User,
				"username": log.Username,
				"uuid": log.Uuid,
			})
		}

		utils.Respond(ctx, template.TemplateLogs.Render(map[string]interface{}{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
			"isPageOne": page == 1,
			"previousPage": page - 1,
			"nextPage": page + 1,
			"logs": formattedLogs,
			"page": page,
		}))
	} else {
		ctx.Redirect(302, "/login")
	}
}

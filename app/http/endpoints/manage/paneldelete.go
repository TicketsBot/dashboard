package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PanelDeleteHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userId := utils.GetUserId(store)

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		messageIdStr := ctx.Param("msg")
		messageId, err := strconv.ParseUint(messageIdStr, 10, 64); if err != nil {
			ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels", guildId))
			return
		}

		// Get object for selected guild
		guild, _ := cache.Instance.GetGuild(guildId, false)

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		// Get CSRF token
		csrfCorrect := ctx.Query("csrf") == store.Get("csrf").(string)
		if !csrfCorrect {
			ctx.Redirect(302, "/")
			return
		}

		go table.DeletePanel(messageId)

		ctx.Redirect(302, fmt.Sprintf("/manage/%d/panels", guildId))
	} else {
		ctx.Redirect(302, "/login")
	}
}

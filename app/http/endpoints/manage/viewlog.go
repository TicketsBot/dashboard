package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func LogViewHandler(ctx *gin.Context) {
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
		guildId, err := strconv.ParseUint(guildIdStr, 10, 64)
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
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		uuid := ctx.Param("uuid")

		// Doesn't need guild = ticket.guild check, since we select where uuid=uuid and guild=guild
		cdnUrl := table.GetCdnUrl(guildId, uuid)

		if cdnUrl == "" {
			ctx.Redirect(302, fmt.Sprintf("/manage/%s/logs/page/1", guild.Id))
			return
		} else {
			req, err := http.NewRequest("GET", cdnUrl, nil); if err != nil {
				ctx.String(500, fmt.Sprintf("Failed to read log: %s", err.Error()))
				return
			}

			client := &http.Client{}
			client.Timeout = 3 * time.Second

			res, err := client.Do(req); if err != nil {
				ctx.String(500, fmt.Sprintf("Failed to read log: %s", err.Error()))
				return
			}
			defer res.Body.Close()

			content, err := ioutil.ReadAll(res.Body); if err != nil {
				ctx.String(500, fmt.Sprintf("Failed to read log: %s", err.Error()))
				return
			}

			ctx.String(200, string(content))
		}
	} else {
		ctx.Redirect(302, "/login")
	}
}

package manage

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func TicketListHandler(ctx *gin.Context) {
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

		// Get object for selected guild
		guild, _ := cache.Instance.GetGuild(guildId, false)

		// Verify the user has permissions to be here
		isAdmin := make(chan bool)
		go utils.IsAdmin(guild, guildId, userId, isAdmin)
		if !<-isAdmin {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		tickets := table.GetOpenTickets(guildId)

		var toFetch []uint64
		for _, ticket := range tickets {
			toFetch = append(toFetch, ticket.Owner)

			for _, idStr := range strings.Split(ticket.Members, ",") {
				if memberId, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					toFetch = append(toFetch, memberId)
				}
			}
		}

		nodes := make(map[uint64]table.UsernameNode)
		for _, node := range table.GetUserNodes(toFetch) {
			nodes[node.Id] = node
		}

		var ticketsFormatted []map[string]interface{}

		for _, ticket := range tickets {
			var membersFormatted []map[string]interface{}
			for index, memberIdStr := range strings.Split(ticket.Members, ",") {
				if memberId, err := strconv.ParseUint(memberIdStr, 10, 64); err == nil {
					if memberId != 0 {
						var separator string
						if index != len(strings.Split(ticket.Members, ",")) - 1 {
							separator = ", "
						}

						membersFormatted = append(membersFormatted, map[string]interface{}{
							"username": nodes[memberId].Name,
							"discrim": nodes[memberId].Discriminator,
							"sep": separator,
						})
					}
				}
			}

			ticketsFormatted = append(ticketsFormatted, map[string]interface{}{
				"uuid": ticket.Uuid,
				"ticketId": ticket.TicketId,
				"username": nodes[ticket.Owner].Name,
				"discrim": nodes[ticket.Owner].Discriminator,
				"members": membersFormatted,
			})
		}

		ctx.HTML(200, "manage/ticketlist", gin.H{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"csrf": store.Get("csrf").(string),
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
			"tickets": ticketsFormatted,
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

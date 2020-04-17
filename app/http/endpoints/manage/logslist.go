package manage

import (
	"context"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

func LogsHandler(ctx *gin.Context) {
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

		pageStr := ctx.Param("page")
		page := 1
		i, err := strconv.Atoi(pageStr)
		if err == nil {
			if i > 0 {
				page = i
			}
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

		pageLimit := 30

		// Get ticket ID from URL
		var ticketId int
		if utils.IsInt(ctx.Query("ticketid")) {
			ticketId, _ = strconv.Atoi(ctx.Query("ticketid"))
		}

		var tickets []table.Ticket

		// Get tickets from DB
		if ticketId > 0 {
			ticketChan := make(chan table.Ticket)
			go table.GetTicketById(guildId, ticketId, ticketChan)
			ticket := <-ticketChan

			if ticket.Uuid != "" && !ticket.IsOpen {
				tickets = append(tickets, ticket)
			}
		} else {
			// make slice of user IDs to filter by
			filteredIds := make([]uint64, 0)

			// Add userid param to slice
			filteredUserId, _ := strconv.ParseUint(ctx.Query("userid"), 10, 64)
			if filteredUserId != 0 {
				filteredIds = append(filteredIds, filteredUserId)
			}

			// Get username from URL
			if username := ctx.Query("username"); username != "" {
				// username -> user id
				rows, err := cache.Instance.PgCache.Query(context.Background(), `select users.user_id from users where "data"->>'Username'=$1 and exists(SELECT FROM members where members.guild_id=$2);`, username, guildId)
				defer rows.Close()
				if err != nil {
					log.Error(err.Error())
					return
				}

				for rows.Next() {
					var filteredId uint64
					if err := rows.Scan(&filteredId); err != nil {
						continue
					}

					if filteredId != 0 {
						filteredIds = append(filteredIds, filteredId)
					}
				}
			}

			if ctx.Query("userid") != "" || ctx.Query("username") != "" {
				tickets = table.GetClosedTicketsByUserId(guildId, filteredIds)
			} else {
				tickets = table.GetClosedTickets(guildId)
			}
		}

		// Select 30 logs + format them
		var formattedLogs []map[string]interface{}
		for i := (page - 1) * pageLimit; i < (page - 1) * pageLimit + pageLimit; i++ {
			if i >= len(tickets) {
				break
			}

			ticket := tickets[i]

			// get username
			user, found := cache.Instance.GetUser(ticket.Owner)
			if !found {
				user, err = rest.GetUser(config.Conf.Bot.Token, ratelimit.Ratelimiter, ticket.Owner)
				if err != nil {
					log.Error(err.Error())
				}
				go cache.Instance.StoreUser(user)
			}

			formattedLogs = append(formattedLogs, map[string]interface{}{
				"ticketid": ticket.TicketId,
				"userid":   ticket.Owner,
				"username": user.Username,
			})
		}

		ctx.HTML(200, "manage/logs",gin.H{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
			"isPageOne": page == 1,
			"previousPage": page - 1,
			"nextPage": page + 1,
			"logs": formattedLogs,
			"page": page,
		})
	} else {
		ctx.Redirect(302, "/login")
	}
}

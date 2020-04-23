package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

const (
	pageLimit = 30
)

func GetLogs(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	page, err := strconv.Atoi(ctx.Param("page"))
	if page < 1 {
		page = 1
	}

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
			rows, err := cache.Instance.PgCache.Query(context.Background(), `select users.user_id from users where LOWER("data"->>'Username') LIKE LOWER($1) and exists(SELECT FROM members where members.guild_id=$2);`, fmt.Sprintf("%%%s%%", username), guildId)
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
	formattedLogs := make([]map[string]interface{}, 0)
	for i := (page - 1) * pageLimit; i < (page-1)*pageLimit+pageLimit; i++ {
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
			"userid":   strconv.FormatUint(ticket.Owner, 10),
			"username": user.Username,
		})
	}

	ctx.JSON(200, formattedLogs)
}

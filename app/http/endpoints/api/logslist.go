package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

const (
	pageLimit = 2
)

func GetLogs(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	before, err := strconv.Atoi(ctx.Query("before"))
	if before < 0 {
		before = 0
	}

	// Get ticket ID from URL
	var ticketId int
	if utils.IsInt(ctx.Query("ticketid")) {
		ticketId, _ = strconv.Atoi(ctx.Query("ticketid"))
	}

	var tickets []database.Ticket

	// Get tickets from DB
	if ticketId > 0 {
		ticket, err := dbclient.Client.Tickets.Get(ticketId, guildId)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
			return
		}

		if ticket.UserId != 0 && !ticket.Open {
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
			tickets, err = dbclient.Client.Tickets.GetMemberClosedTickets(guildId, filteredIds, pageLimit, before)
		} else {
			tickets, err = dbclient.Client.Tickets.GetGuildClosedTickets(guildId, pageLimit, before)
		}

		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error": err.Error(),
			})
			return
		}
	}

	// Select 30 logs + format them
	formattedLogs := make([]map[string]interface{}, 0)
	for _, ticket := range tickets {
		// get username
		user, found := cache.Instance.GetUser(ticket.UserId)
		if !found {
			user, err = rest.GetUser(config.Conf.Bot.Token, ratelimit.Ratelimiter, ticket.UserId)
			if err != nil {
				log.Error(err.Error())
			}
			go cache.Instance.StoreUser(user)
		}

		formattedLogs = append(formattedLogs, map[string]interface{}{
			"ticketid": ticket.Id,
			"userid":   strconv.FormatUint(ticket.UserId, 10),
			"username": user.Username,
		})
	}

	ctx.JSON(200, formattedLogs)
}

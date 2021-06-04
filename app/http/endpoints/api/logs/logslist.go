package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"net/http"
	"strconv"
)

const (
	pageLimit = 30
)

type filterType uint8

const (
	filterTypeNone filterType = iota
	filterTypeTicketId
	filterTypeUsername
	filterTypeUserId
)

func GetLogs(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	before, _ := strconv.Atoi(ctx.Query("before"))
	if before < 0 {
		before = 0
	}

	var tickets []database.Ticket
	var status int

	filterType := getFilterType(ctx)
	switch filterType {
	case filterTypeNone:
		tickets, err = getTickets(guildId, before)
	case filterTypeTicketId:
		tickets, status,
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
				"error":   err.Error(),
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
				"error":   err.Error(),
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
			user, err = rest.GetUser(botContext.Token, botContext.RateLimiter, ticket.UserId)
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

func getFilterType(ctx *gin.Context) filterType {
	if ctx.Query("ticket_id") != "" {
		return filterTypeTicketId
	} else if ctx.Query("username") != "" {
		return filterTypeUsername
	} else if ctx.Query("userid") != "" {
		return filterTypeUserId
	} else {
		return filterTypeNone
	}
}

func getTickets(guildId uint64, before int) ([]database.Ticket, int, error) {
	tickets, err := dbclient.Client.Tickets.GetGuildClosedTickets(guildId, pageLimit, before)
}

// (tickets, statusCode, error)
func getTicketsByTicketId(guildId uint64, ctx *gin.Context) ([]database.Ticket, int, error) {
	ticketId, err := strconv.Atoi(ctx.Query("ticketid"))
	if err != nil {
		return nil, 400, err
	}

	ticket, err := dbclient.Client.Tickets.Get(ticketId, guildId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if ticket.Id == 0 {
		return nil, http.StatusNotFound, errors.New("ticket not found")
	}

	return []database.Ticket{ticket}, http.StatusOK, nil
}

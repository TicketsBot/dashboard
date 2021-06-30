package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"math"
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

type transcript struct {
	TicketId    int     `json:"ticket_id"`
	Username    string  `json:"username"`
	CloseReason *string `json:"close_reason"`
}

func ListTranscripts(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// db functions will handle if 0
	before, _ := strconv.Atoi(ctx.Query("before"))
	after, _ := strconv.Atoi(ctx.Query("after"))

	var tickets []database.TicketWithCloseReason
	var status int

	filterType := getFilterType(ctx)
	switch filterType {
	case filterTypeNone:
		tickets, status, err = getTickets(guildId, before, after)
	case filterTypeTicketId:
		tickets, status, err = getTicketsByTicketId(guildId, ctx)
	case filterTypeUsername:
		tickets, status, err = getTicketsByUsername(guildId, before, after, ctx)
	case filterTypeUserId:
		tickets, status, err = getTicketsByUserId(guildId, before, after, ctx)
	}

	if err != nil {
		ctx.JSON(status, utils.ErrorJson(err))
		return
	}

	// Create a mapping user_id -> username so we can skip duplicates
	usernames := make(map[uint64]string)
	for _, ticket := range tickets {
		if _, ok := usernames[ticket.UserId]; ok {
			continue // don't fetch again
		}

		// check cache, for some reason botContext.GetUser doesn't do this
		user, ok := cache.Instance.GetUser(ticket.UserId)
		if ok {
			usernames[ticket.UserId] = user.Username
		} else {
			user, err = botContext.GetUser(ticket.UserId)
			if err != nil { // TODO: Log
				usernames[ticket.UserId] = "Unknown User"
			} else {
				usernames[ticket.UserId] = user.Username
			}
		}
	}

	transcripts := make([]transcript, len(tickets))
	for i, ticket := range tickets {
		transcripts[i] = transcript{
			TicketId:    ticket.Id,
			Username:    usernames[ticket.UserId],
			CloseReason: ticket.CloseReason,
		}
	}

	ctx.JSON(200, transcripts)
}

func getFilterType(ctx *gin.Context) filterType {
	if ctx.Query("ticketid") != "" {
		return filterTypeTicketId
	} else if ctx.Query("username") != "" {
		return filterTypeUsername
	} else if ctx.Query("userid") != "" {
		return filterTypeUserId
	} else {
		return filterTypeNone
	}
}

func getTickets(guildId uint64, before, after int) ([]database.TicketWithCloseReason, int, error) {
	var tickets []database.TicketWithCloseReason
	var err error

	if before <= 0 && after <= 0 {
		tickets, err = dbclient.Client.Tickets.GetGuildClosedTicketsBeforeWithCloseReason(guildId, pageLimit, math.MaxInt32)
	} else if before > 0 {
		tickets, err = dbclient.Client.Tickets.GetGuildClosedTicketsBeforeWithCloseReason(guildId, pageLimit, before)
	} else { // after > 0
		// returns in ascending order, must reverse
		tickets, err = dbclient.Client.Tickets.GetGuildClosedTicketsAfterWithCloseReason(guildId, pageLimit, after)
		if err == nil {
			reverse(tickets)
		}
	}


	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}

	return tickets, status, err
}

// (tickets, statusCode, error)
func getTicketsByTicketId(guildId uint64, ctx *gin.Context) ([]database.TicketWithCloseReason, int, error) {
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

	closeReason, ok, err := dbclient.Client.CloseReason.Get(guildId, ticketId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	data := database.TicketWithCloseReason{
		Ticket: ticket,
	}

	if ok {
		data.CloseReason = &closeReason
	}

	return []database.TicketWithCloseReason{data}, http.StatusOK, nil
}

// (tickets, statusCode, error)
func getTicketsByUsername(guildId uint64, before, after int, ctx *gin.Context) ([]database.TicketWithCloseReason, int, error) {
	username := ctx.Query("username")

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	members, err := botContext.SearchMembers(guildId, username)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	userIds := make([]uint64, len(members)) // capped at 100
	for i, member := range members {
		userIds[i] = member.User.Id
	}

	var tickets []database.TicketWithCloseReason
	if before <= 0 && after <= 0 {
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyBeforeWithCloseReason(guildId, userIds, math.MaxInt32, pageLimit)
	} else if before > 0 {
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyBeforeWithCloseReason(guildId, userIds, before, pageLimit)
	} else { // after > 0
		// returns in ascending order, must reverse
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyAfterWithCloseReason(guildId, userIds, after, pageLimit)
		if err == nil {
			reverse(tickets)
		}
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return tickets, http.StatusOK, nil
}

// (tickets, statusCode, error)
func getTicketsByUserId(guildId uint64, before, after int, ctx *gin.Context) ([]database.TicketWithCloseReason, int, error) {
	userId, err := strconv.ParseUint(ctx.Query("userid"), 10, 64)
	if err != nil {
		return nil, 400, err
	}

	var tickets []database.TicketWithCloseReason
	if before <= 0 && after <= 0 {
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyBeforeWithCloseReason(guildId, []uint64{userId}, math.MaxInt32, pageLimit)
	} else if before > 0 {
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyBeforeWithCloseReason(guildId, []uint64{userId}, before, pageLimit)
	} else { // after > 0
		// returns in ascending order, must reverse
		tickets, err = dbclient.Client.Tickets.GetClosedByAnyAfterWithCloseReason(guildId, []uint64{userId}, after, pageLimit)
		if err == nil {
			reverse(tickets)
		}
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return tickets, http.StatusOK, nil
}

func reverse(slice []database.TicketWithCloseReason) {
	if len(slice) == 0 {
		return
	}

	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

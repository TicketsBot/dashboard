package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
)

const pageLimit = 15

type transcriptMetadata struct {
	TicketId      int     `json:"ticket_id"`
	Username      string  `json:"username"`
	CloseReason   *string `json:"close_reason"`
	ClosedBy      *uint64 `json:"closed_by"`
	Rating        *uint8  `json:"rating"`
	HasTranscript bool    `json:"has_transcript"`
}

func ListTranscripts(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var queryOptions wrappedQueryOptions
	if err := ctx.BindJSON(&queryOptions); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	opts, err := queryOptions.toQueryOptions(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	tickets, err := dbclient.Client.Tickets.GetByOptions(opts)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
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

	// Get ratings
	ticketIds := make([]int, len(tickets))
	for i, ticket := range tickets {
		ticketIds[i] = ticket.Id
	}

	ratings, err := dbclient.Client.ServiceRatings.GetMulti(guildId, ticketIds)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Get close reasons
	closeReasons, err := dbclient.Client.CloseReason.GetMulti(guildId, ticketIds)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	transcripts := make([]transcriptMetadata, len(tickets))
	for i, ticket := range tickets {
		transcript := transcriptMetadata{
			TicketId:      ticket.Id,
			Username:      usernames[ticket.UserId],
			HasTranscript: ticket.HasTranscript,
		}

		if v, ok := ratings[ticket.Id]; ok {
			transcript.Rating = &v
		}

		if v, ok := closeReasons[ticket.Id]; ok {
			transcript.CloseReason = v.Reason
			transcript.ClosedBy = v.ClosedBy
		}

		transcripts[i] = transcript
	}

	ctx.JSON(200, transcripts)
}

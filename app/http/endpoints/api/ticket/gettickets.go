package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
	"golang.org/x/sync/errgroup"
)

func GetTickets(ctx *gin.Context) {
	type WithUser struct {
		TicketId int        `json:"id"`
		User     *user.User `json:"user,omitempty"`
	}

	guildId := ctx.Keys["guildid"].(uint64)

	tickets, err := database.Client.Tickets.GetGuildOpenTickets(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data := make([]WithUser, len(tickets))

	group, _ := errgroup.WithContext(context.Background())

	for i, ticket := range tickets {
		i := i
		ticket := ticket

		group.Go(func() error {
			user, ok := cache.Instance.GetUser(ticket.UserId)

			data[i] = WithUser{
				TicketId: ticket.Id,
			}

			if ok {
				data[i].User = &user
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, data)
}

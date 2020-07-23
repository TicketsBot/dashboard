package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func GetTickets(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	tickets, err := database.Client.Tickets.GetGuildOpenTickets(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ticketsFormatted := make([]map[string]interface{}, len(tickets))

	group, _ := errgroup.WithContext(context.Background())

	for i, ticket := range tickets {
		i := i
		ticket := ticket

		group.Go(func() error {
			members, err := database.Client.TicketMembers.Get(guildId, ticket.Id)
			if err != nil {
				return err
			}

			membersFormatted := make([]map[string]interface{}, 0)
			for _, userId := range members {
				user, _ := cache.Instance.GetUser(userId)

				membersFormatted = append(membersFormatted, map[string]interface{}{
					"id": strconv.FormatUint(userId, 10),
					"username": user.Username,
					"discrim":  fmt.Sprintf("%04d", user.Discriminator),
				})
			}

			owner, _ := cache.Instance.GetUser(ticket.UserId)

			ticketsFormatted[len(tickets) - 1 - i] =  map[string]interface{}{
				"ticketId": ticket.Id,
				"username": owner.Username,
				"discrim":  fmt.Sprintf("%04d", owner.Discriminator),
				"members":  membersFormatted,
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, ticketsFormatted)
}

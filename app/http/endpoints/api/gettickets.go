package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetTickets(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	tickets := table.GetOpenTickets(guildId)
	ticketsFormatted := make([]map[string]interface{}, 0)

	for _, ticket := range tickets {
		membersFormatted := make([]map[string]interface{}, 0)
		for index, memberIdStr := range strings.Split(ticket.Members, ",") {
			if memberId, err := strconv.ParseUint(memberIdStr, 10, 64); err == nil {
				if memberId != 0 {
					var separator string
					if index != len(strings.Split(ticket.Members, ","))-1 {
						separator = ", "
					}

					member, _ := cache.Instance.GetUser(memberId)
					membersFormatted = append(membersFormatted, map[string]interface{}{
						"username": member.Username,
						"discrim":  fmt.Sprintf("%04d", member.Discriminator),
						"sep":      separator,
					})
				}
			}
		}

		owner, _ := cache.Instance.GetUser(ticket.Owner)
		ticketsFormatted = append(ticketsFormatted, map[string]interface{}{
			"uuid":     ticket.Uuid,
			"ticketId": ticket.TicketId,
			"username": owner.Username,
			"discrim":  fmt.Sprintf("%04d", owner.Discriminator),
			"members":  membersFormatted,
		})
	}

	ctx.JSON(200, ticketsFormatted)
}

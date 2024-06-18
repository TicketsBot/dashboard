package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
)

type ticketResponse struct {
	TicketId   int        `json:"id"`
	PanelTitle string     `json:"panel_title"`
	User       *user.User `json:"user,omitempty"`
}

func GetTickets(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	tickets, err := database.Client.Tickets.GetGuildOpenTickets(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panels, err := database.Client.Panel.GetByGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panelTitles := make(map[int]string)
	for _, panel := range panels {
		panelTitles[panel.PanelId] = panel.Title
	}

	// Get user objects
	userIds := make([]uint64, len(tickets))
	for i, ticket := range tickets {
		userIds[i] = ticket.UserId
	}

	users, err := cache.Instance.GetUsers(context.Background(), userIds)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data := make([]ticketResponse, len(tickets))
	for i, ticket := range tickets {
		var user *user.User
		if tmp, ok := users[ticket.UserId]; ok {
			user = &tmp
		}

		panelTitle := "Unknown"
		if ticket.PanelId != nil {
			if tmp, ok := panelTitles[*ticket.PanelId]; ok {
				panelTitle = tmp
			}
		}

		data[i] = ticketResponse{
			TicketId:   ticket.Id,
			PanelTitle: panelTitle,
			User:       user,
		}
	}

	ctx.JSON(200, data)
}

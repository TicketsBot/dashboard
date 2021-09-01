package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/utils"
)

type wrappedQueryOptions struct {
	Id       int    `json:"id,string"`
	Username string `json:"username"`
	UserId   uint64 `json:"user_id,string"`
	PanelId  int    `json:"panel_id"`
	Page     int    `json:"page"`
}

func (o *wrappedQueryOptions) toQueryOptions(guildId uint64) (database.TicketQueryOptions, error) {
	var userIds []uint64
	if len(o.Username) > 0 {
		var err error
		userIds, err = usernameToIds(guildId, o.Username)
		if err != nil {
			return database.TicketQueryOptions{}, err
		}

		// TODO: Do this better
		if len(userIds) == 0 {
			return database.TicketQueryOptions{}, errors.New("User not found")
		}
	}

	if o.UserId != 0 {
		userIds = append(userIds, o.UserId)
	}

	var offset int
	if o.Page > 1 {
		offset = pageLimit * (o.Page - 1)
	}

	opts := database.TicketQueryOptions{
		Id:      o.Id,
		GuildId: guildId,
		UserIds: userIds,
		Open:    utils.BoolPtr(false),
		PanelId: o.PanelId,
		Order:   database.OrderTypeDescending,
		Limit:   pageLimit,
		Offset:  offset,
	}
	return opts, nil
}

func usernameToIds(guildId uint64, username string) ([]uint64, error) {
	if len(username) > 32 {
		return nil, errors.New("username too long")
	}

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return nil, err
	}

	members, err := botContext.SearchMembers(guildId, username)
	if err != nil {
		return nil, err
	}

	userIds := make([]uint64, len(members)) // capped at 100
	for i, member := range members {
		userIds[i] = member.User.Id
	}

	return userIds, nil
}

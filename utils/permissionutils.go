package utils

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

func IsAdmin(g guild.Guild, guildId, userId uint64, res chan bool) {
	if Contains(config.Conf.Admins, strconv.Itoa(int(userId))) {
		res <- true
	}

	if g.Owner {
		res <- true
	}

	if table.IsAdmin(guildId, userId) {
		res <- true
	}

	if g.Permissions & 0x8 != 0 {
		res <- true
	}

	adminRolesChan := make(chan []uint64)
	go table.GetAdminRoles(guildId, adminRolesChan)
	adminRoles := <- adminRolesChan

	userRoles, found := getRoles(guildId, userId)

	hasAdminRole := false
	if found {
		for _, userRole := range userRoles {
			for _, adminRole := range adminRoles {
				if userRole == adminRole {
					hasAdminRole = true
					break
				}
			}
		}
	}

	if hasAdminRole {
		res <- true
	}

	res <- false
}

func getRoles(guildId, userId uint64) ([]uint64, bool) {
	member, found := cache.Instance.GetMember(guildId, userId)
	if !found { // get from rest
		var err error
		member, err = rest.GetGuildMember(config.Conf.Bot.Token, ratelimit.Ratelimiter, guildId, userId)
		if err != nil {
			return nil, false
		}

		// cache
		cache.Instance.StoreMember(member, guildId)
	}

	return member.Roles, true
}

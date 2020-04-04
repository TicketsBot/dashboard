package utils

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/guild"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/robfig/go-cache"
	"strconv"
	"time"
)

var roleCache = cache.New(time.Minute, time.Minute)

func IsAdmin(guild objects.Guild, guildId, userId uint64, res chan bool) {
	if Contains(config.Conf.Admins, strconv.Itoa(int(userId))) {
		res <- true
	}

	if guild.Owner {
		res <- true
	}

	if table.IsAdmin(guildId, userId) {
		res <- true
	}

	if guild.Permissions & 0x8 != 0 {
		res <- true
	}

	userRolesChan := make(chan []uint64)
	go table.GetCachedRoles(guildId, userId, userRolesChan)
	userRoles := <-userRolesChan

	adminRolesChan := make(chan []uint64)
	go table.GetAdminRoles(guildId, adminRolesChan)
	adminRoles := <- adminRolesChan

	hasAdminRole := false
	for _, userRole := range userRoles {
		for _, adminRole := range adminRoles {
			if userRole == adminRole {
				hasAdminRole = true
				break
			}
		}
	}

	if hasAdminRole {
		res <- true
	}

	res <- false
}

func GetRolesRest(store sessions.Session, guildId, userId uint64) *[]uint64 {
	var member objects.Member
	endpoint := guild.GetGuildMember(guildId, userId)

	if err, _ := endpoint.Request(store, nil, nil, &member); err != nil {
		log.Error(err.Error())
		return nil
	}

	roles := []uint64(member.Roles)
	return &roles
}

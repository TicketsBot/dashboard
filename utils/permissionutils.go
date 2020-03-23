package utils

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/guild"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"strconv"
)

func IsAdmin(guild objects.Guild, guildId, user int64, res chan bool) {
	userIdStr := strconv.Itoa(int(user))

	if Contains(config.Conf.Admins, userIdStr) {
		res <- true
	}

	if guild.Owner {
		res <- true
	}

	if table.IsAdmin(guildId, user) {
		res <- true
	}

	res <- false
}

func GetRolesRest(store sessions.Session, guildId, userId int64, res chan *[]int64) {
	var member objects.Member
	endpoint := guild.GetGuildMember(int(guildId), int(userId))

	if err, _ := endpoint.Request(store, nil, nil, &member); err != nil {
		res <- nil
	}

	roles := make([]int64, 0)
	for _, role := range member.Roles {
		int, err := strconv.ParseInt(role, 10, 64); if err != nil {
			log.Error(err.Error())
			continue
		}

		roles = append(roles, int)
	}

	res <- &roles
}

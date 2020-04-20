package utils

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/apex/log"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/permission"
	"github.com/rxdn/gdl/rest"
	"golang.org/x/sync/errgroup"
	"strconv"
	"sync"
)

func IsAdmin(g guild.Guild, userId uint64, res chan bool) {
	if Contains(config.Conf.Admins, strconv.Itoa(int(userId))) {
		res <- true
	}

	if g.OwnerId == userId {
		res <- true
	}

	if table.IsAdmin(g.Id, userId) {
		res <- true
	}

	userRoles, _ := getRoles(g.Id, userId)

	// check if user has administrator permission
	if hasAdministratorPermission(g.Id, userRoles) {
		res <- true
	}

	adminRolesChan := make(chan []uint64)
	go table.GetAdminRoles(g.Id, adminRolesChan)
	adminRoles := <-adminRolesChan

	hasTicketAdminRole := false
	for _, userRole := range userRoles {
		for _, adminRole := range adminRoles {
			if userRole == adminRole {
				hasTicketAdminRole = true
				break
			}
		}
	}

	if hasTicketAdminRole {
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

func hasAdministratorPermission(guildId uint64, roles []uint64) bool {
	var lock sync.Mutex
	var hasAdministrator bool

	group, _ := errgroup.WithContext(context.Background())
	for _, roleId := range roles {
		group.Go(func() error {
			roleHasAdmin, err := roleHasAdministrator(guildId, roleId)
			if err != nil {
				return err
			}

			if roleHasAdmin {
				lock.Lock()
				hasAdministrator = true
				lock.Unlock()
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return false
	}

	return hasAdministrator
}

func roleHasAdministrator(guildId, roleId uint64) (bool, error) {
	role, found := cache.Instance.GetRole(roleId)
	if !found {
		roles, err := rest.GetGuildRoles(config.Conf.Bot.Token, ratelimit.Ratelimiter, guildId)
		if err != nil {
			log.Error(err.Error())
			return false, err
		}

		go cache.Instance.StoreRoles(roles, guildId)
		for _, r := range roles {
			if r.Id == roleId {
				role = r
				found = true
				break
			}
		}
		if !found {
			return false, errors.New("role does not exist")
		}
	}

	return permission.HasPermissionRaw(role.Permissions, permission.Administrator), nil
}

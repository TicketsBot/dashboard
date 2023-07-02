package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	syncutils "github.com/TicketsBot/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
	"golang.org/x/sync/errgroup"
	"sort"
	"strconv"
)

func GetMembers(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	teamId := ctx.Param("teamid")
	if teamId == "default" {
		getDefaultMembers(ctx, guildId)
	} else {
		parsed, err := strconv.Atoi(teamId)
		if err != nil {
			ctx.JSON(400, utils.ErrorStr("Invalid team ID"))
			return
		}

		getTeamMembers(ctx, parsed, guildId)
	}
}

func getDefaultMembers(ctx *gin.Context, guildId uint64) {
	group, _ := errgroup.WithContext(context.Background())

	// get IDs of support users & roles
	var userIds []uint64
	group.Go(func() (err error) {
		userIds, err = dbclient.Client.Permissions.GetSupport(guildId)
		return
	})

	var roleIds []uint64
	group.Go(func() (err error) {
		roleIds, err = dbclient.Client.RolePermissions.GetSupportRoles(guildId)
		return
	})

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data, err := formatMembers(guildId, userIds, roleIds)
	if err == nil {
		ctx.JSON(200, data)
	} else {
		ctx.JSON(500, utils.ErrorJson(err))
	}
}

func getTeamMembers(ctx *gin.Context, teamId int, guildId uint64) {
	// Verify team exists
	exists, err := dbclient.Client.SupportTeam.Exists(teamId, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !exists {
		ctx.JSON(404, utils.ErrorStr("Support team with provided ID not found"))
		return
	}

	group, _ := errgroup.WithContext(context.Background())

	// get IDs of support users & roles
	var userIds []uint64
	group.Go(func() (err error) {
		userIds, err = dbclient.Client.SupportTeamMembers.Get(teamId)
		return
	})

	var roleIds []uint64
	group.Go(func() (err error) {
		roleIds, err = dbclient.Client.SupportTeamRoles.Get(teamId)
		return
	})

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data, err := formatMembers(guildId, userIds, roleIds)
	if err == nil {
		ctx.JSON(200, data)
	} else {
		ctx.JSON(500, utils.ErrorJson(err))
	}
}

func formatMembers(guildId uint64, userIds, roleIds []uint64) ([]entity, error) {
	ctx, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return nil, err
	}

	// map role ids to names
	data := make([]entity, 0)
	for _, roleId := range roleIds {
		data = append(data, entity{
			Id:   roleId,
			Type: entityTypeRole,
		})
	}

	// map user ids to names & discrims
	group, _ := errgroup.WithContext(context.Background())

	users := make(chan user.User)
	wg := syncutils.NewChannelWaitGroup()
	wg.Add(len(userIds))

	for _, userId := range userIds {
		userId := userId

		group.Go(func() error {
			defer wg.Done()

			user, err := ctx.GetUser(userId)
			if err != nil {
				// TODO: Log w sentry
				return nil // We should skip the error, since it's probably 403 / 404 etc
			}

			users <- user
			return nil
		})
	}

	group.Go(func() error {
	loop:
		for {
			select {
			case <-wg.Wait():
				break loop
			case user := <-users:
				data = append(data, entity{
					Id:   user.Id,
					Name: fmt.Sprintf("%s", user.Username),
					Type: entityTypeUser,
				})
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	// sort
	sort.Slice(data, func(i, j int) bool {
		if data[i].Type == data[j].Type {
			return data[i].Id < data[j].Id
		} else {
			return data[i].Type > data[j].Type
		}
	})

	return data, nil
}

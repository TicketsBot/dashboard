package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

func RemoveMember(ctx *gin.Context) {
	guildId, selfId := ctx.Keys["guildid"].(uint64), ctx.Keys["userid"].(uint64)

	snowflake, err := strconv.ParseUint(ctx.Param("snowflake"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get entity type
	typeParsed, err := strconv.Atoi(ctx.Query("type"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	entityType, ok := entityTypes[typeParsed]
	if !ok {
		ctx.JSON(400, utils.ErrorStr("Invalid entity type"))
		return
	}

	teamId := ctx.Param("teamid")
	if teamId == "default" {
		removeDefaultMember(ctx, guildId, selfId, snowflake, entityType)
	} else {
		parsed, err := strconv.Atoi(teamId)
		if err != nil {
			ctx.JSON(400, utils.ErrorStr("Invalid team ID"))
			return
		}

		removeTeamMember(ctx, parsed, guildId, snowflake, entityType)
	}
}

func removeDefaultMember(ctx *gin.Context, guildId, selfId, snowflake uint64, entityType entityType) {
	// permission check
	var isAdmin bool
	var err error
	switch entityType {
	case entityTypeUser:
		isAdmin, err = dbclient.Client.Permissions.IsAdmin(guildId, snowflake)
	case entityTypeRole:
		isAdmin, err = dbclient.Client.RolePermissions.IsAdmin(snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// only guild owner can remove admins
	if isAdmin {
		botCtx, err := botcontext.ContextForGuild(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// TODO: Use proper context
		guild, err := botCtx.GetGuild(context.Background(), guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if guild.OwnerId != selfId {
			ctx.JSON(403, utils.ErrorStr("Only the server owner can remove admins"))
			return
		}
	}

	switch entityType {
	case entityTypeUser:
		err = dbclient.Client.Permissions.RemoveSupport(guildId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.RolePermissions.RemoveSupport(guildId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Remove on-call role
	metadata, err := dbclient.Client.GuildMetadata.Get(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if metadata.OnCallRole != nil {
		botContext, err := botcontext.ContextForGuild(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if entityType == entityTypeUser {
			// If the member  is not in the guild we do not have to worry
			// TODO: Use proper context
			member, err := botContext.GetGuildMember(context.Background(), guildId, snowflake)
			if err == nil {
				if member.HasRole(*metadata.OnCallRole) {
					// Attempt to remove role but ignore failure
					// TODO: Use proper context
					_ = botContext.RemoveGuildMemberRole(context.Background(), guildId, snowflake, *metadata.OnCallRole)
				}
			} else {
				if err, ok := err.(request.RestError); !ok || err.StatusCode != 404 {
					ctx.JSON(500, utils.ErrorJson(err))
					return
				}
			}
		} else if entityType == entityTypeRole {
			// Recreate role
			if err := dbclient.Client.GuildMetadata.SetOnCallRole(guildId, nil); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			// TODO: Use proper context
			if err := botContext.DeleteGuildRole(context.Background(), guildId, *metadata.OnCallRole); err != nil && !isUnknownRoleError(err) {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			if _, err := createOnCallRole(botContext, guildId, nil); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		} else {
			ctx.JSON(500, utils.ErrorStr("Infallible"))
			return
		}
	}

	ctx.JSON(200, utils.SuccessResponse)
}

func removeTeamMember(ctx *gin.Context, teamId int, guildId, snowflake uint64, entityType entityType) {
	team, exists, err := dbclient.Client.SupportTeam.GetById(guildId, teamId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !exists {
		ctx.JSON(404, utils.ErrorStr("Support team with provided ID not found"))
		return
	}

	// Remove from DB
	switch entityType {
	case entityTypeUser:
		err = dbclient.Client.SupportTeamMembers.Delete(teamId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.SupportTeamRoles.Delete(teamId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Remove on-call role
	if team.OnCallRole != nil {
		botContext, err := botcontext.ContextForGuild(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if entityType == entityTypeUser {
			// If the member  is not in the guild we do not have to worry
			// TODO: Use proper context
			member, err := botContext.GetGuildMember(context.Background(), guildId, snowflake)
			if err == nil {
				if member.HasRole(*team.OnCallRole) {
					// Attempt to remove role but ignore failure
					// TODO: Use proper context
					_ = botContext.RemoveGuildMemberRole(context.Background(), guildId, snowflake, *team.OnCallRole)
				}
			} else {
				var err request.RestError
				if !errors.As(err, &err) || err.StatusCode != 404 {
					ctx.JSON(500, utils.ErrorJson(err))
					return
				}
			}

			// TODO: Use proper context
			_ = botContext.RemoveGuildMemberRole(context.Background(), guildId, snowflake, *team.OnCallRole)
		} else if entityType == entityTypeRole {
			// Recreate role
			if err := dbclient.Client.SupportTeam.SetOnCallRole(teamId, nil); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			// TODO: Use proper context
			if err := botContext.DeleteGuildRole(context.Background(), guildId, *team.OnCallRole); err != nil && !isUnknownRoleError(err) {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			if _, err := createOnCallRole(botContext, guildId, &team); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		} else {
			ctx.JSON(500, utils.ErrorStr("Infallible"))
		}
	}

	ctx.JSON(200, utils.SuccessResponse)
}

func createOnCallRole(botContext *botcontext.BotContext, guildId uint64, team *database.SupportTeam) (uint64, error) {
	var roleName string
	if team == nil {
		roleName = "On Call" // TODO: Translate
	} else {
		roleName = utils.StringMax(fmt.Sprintf("On Call - %s", team.Name), 100)
	}

	data := rest.GuildRoleData{
		Name:        roleName,
		Hoist:       utils.Ptr(false),
		Mentionable: utils.Ptr(false),
	}

	// TODO: Use proper context
	role, err := botContext.CreateGuildRole(context.Background(), guildId, data)
	if err != nil {
		return 0, err
	}

	if team == nil {
		if err := dbclient.Client.GuildMetadata.SetOnCallRole(guildId, &role.Id); err != nil {
			return 0, err
		}
	} else {
		if err := dbclient.Client.SupportTeam.SetOnCallRole(team.Id, &role.Id); err != nil {
			return 0, err
		}
	}

	return role.Id, nil
}

func isUnknownRoleError(err error) bool {
	var restErr request.RestError
	if errors.As(err, &restErr) && restErr.ApiError.Message == "Unknown Role" {
		return true
	}

	return false
}

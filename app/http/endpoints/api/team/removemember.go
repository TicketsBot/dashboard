package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
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

		guild, err := botCtx.GetGuild(guildId)
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

	ctx.JSON(200, utils.SuccessResponse)
}

func removeTeamMember(ctx *gin.Context, teamId int, guildId, snowflake uint64, entityType entityType) {
	exists, err := dbclient.Client.SupportTeam.Exists(teamId, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !exists {
		ctx.JSON(404, utils.ErrorStr("Support team with provided ID not found"))
		return
	}

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

	ctx.JSON(200, utils.SuccessResponse)
}
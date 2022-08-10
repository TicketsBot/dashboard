package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddMember(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

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

	if entityType == entityTypeRole && snowflake == guildId {
		ctx.JSON(400, utils.ErrorStr("You cannot add the @everyone role as staff"))
		return
	}

	teamId := ctx.Param("teamid")
	if teamId == "default" {
		addDefaultMember(ctx, guildId, snowflake, entityType)
	} else {
		parsed, err := strconv.Atoi(teamId)
		if err != nil {
			ctx.JSON(400, utils.ErrorStr("Invalid team ID"))
			return
		}

		addTeamMember(ctx, parsed, guildId, snowflake, entityType)
	}
}

func addDefaultMember(ctx *gin.Context, guildId, snowflake uint64, entityType entityType) {
	var err error
	switch entityType {
	case entityTypeUser:
		err = dbclient.Client.Permissions.AddSupport(guildId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.RolePermissions.AddSupport(guildId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

func addTeamMember(ctx *gin.Context, teamId int, guildId, snowflake uint64, entityType entityType) {
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
		err = dbclient.Client.SupportTeamMembers.Add(teamId, snowflake)
	case entityTypeRole:
		err = dbclient.Client.SupportTeamRoles.Add(teamId, snowflake)
	}

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

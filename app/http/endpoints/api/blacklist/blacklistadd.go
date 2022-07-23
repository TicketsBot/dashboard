package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	id, err := strconv.ParseUint(ctx.Param("user"), 10, 64)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// Max of 250 blacklisted users
	count, err := database.Client.Blacklist.GetBlacklistedCount(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if count >= 250 {
		ctx.JSON(400, utils.ErrorStr("Blacklist limit (250) reached: consider using a role instead"))
		return
	}

	permLevel, err := utils.GetPermissionLevel(guildId, id)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if permLevel > permission.Everyone {
		ctx.JSON(400, utils.ErrorStr("You cannot blacklist staff members!"))
		return
	}

	if err = database.Client.Blacklist.Add(guildId, id); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

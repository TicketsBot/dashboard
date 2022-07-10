package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func DeleteIntegrationHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	integrationId, err := strconv.Atoi(ctx.Param("integrationid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid integration ID"))
		return
	}

	integration, ok, err := dbclient.Client.CustomIntegrations.Get(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !ok {
		ctx.JSON(404, utils.ErrorStr("Integration not found"))
		return
	}

	// Check if the user has permission to manage this integration
	if integration.OwnerId != userId {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to delete this integration"))
		return
	}

	if err := dbclient.Client.CustomIntegrations.Delete(integration.Id); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

type detailedResponse struct {
	database.CustomIntegration
	Placeholders []database.CustomIntegrationPlaceholder `json:"placeholders"`
	Headers      []database.CustomIntegrationHeader      `json:"headers"`
	Secrets      []database.CustomIntegrationSecret      `json:"secrets"`
}

func GetIntegrationDetailedHandler(ctx *gin.Context) {
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

	// Check if the user has permission to view this integration
	if integration.OwnerId != userId {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to view this integration"))
		return
	}

	// Get placeholders
	placeholders, err := dbclient.Client.CustomIntegrationPlaceholders.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Don't serve null
	if placeholders == nil {
		placeholders = make([]database.CustomIntegrationPlaceholder, 0)
	}

	// Get headers
	headers, err := dbclient.Client.CustomIntegrationHeaders.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Don't serve null
	if headers == nil {
		headers = make([]database.CustomIntegrationHeader, 0)
	}

	// Get secrets
	secrets, err := dbclient.Client.CustomIntegrationSecrets.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Don't serve null
	if secrets == nil {
		secrets = make([]database.CustomIntegrationSecret, 0)
	}

	ctx.JSON(200, detailedResponse{
		CustomIntegration: integration,
		Placeholders:      placeholders,
		Headers:           headers,
		Secrets:           secrets,
	})
}

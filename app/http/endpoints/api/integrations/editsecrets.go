package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func UpdateIntegrationSecretsHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	integrationId, err := strconv.Atoi(ctx.Param("integrationid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid integration ID"))
		return
	}

	// Check integration is active
	active, err := dbclient.Client.CustomIntegrationGuilds.IsActive(integrationId, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !active {
		ctx.JSON(400, utils.ErrorStr("Integration is not active"))
		return
	}

	var data activateIntegrationBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}
	// Check the secret values are valid
	secrets, err := dbclient.Client.CustomIntegrationSecrets.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if len(secrets) != len(data.Secrets) {
		ctx.JSON(400, utils.ErrorStr("Invalid secret values"))
		return
	}

	// Since we've checked the length, we can just iterate over the secrets and they're guaranteed to be correct
	secretMap := make(map[int]string)
	for secretName, value := range data.Secrets {
		if len(value) == 0 || len(value) > 255 {
			ctx.JSON(400, utils.ErrorStr("Secret values must be between 1 and 255 characters"))
			return
		}

		found := false

	inner:
		for _, secret := range secrets {
			if secret.Name == secretName {
				found = true
				secretMap[secret.Id] = value
				break inner
			}
		}

		if !found {
			ctx.JSON(400, utils.ErrorStr("Invalid secret values"))
			return
		}
	}

	if err := dbclient.Client.CustomIntegrationSecretValues.UpdateAll(guildId, integrationId, secretMap); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

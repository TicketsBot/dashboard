package api

import (
	"encoding/json"
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type activateIntegrationBody struct {
	Secrets map[string]string `json:"secrets"`
}

func ActivateIntegrationHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)
	guildId := ctx.Keys["guildid"].(uint64)

	activeCount, err := dbclient.Client.CustomIntegrationGuilds.GetGuildIntegrationCount(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if activeCount >= 5 {
		ctx.JSON(400, utils.ErrorStr("You can only have 5 integrations active at once"))
		return
	}

	integrationId, err := strconv.Atoi(ctx.Param("integrationid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid integration ID"))
		return
	}

	var data activateIntegrationBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	integration, ok, err := dbclient.Client.CustomIntegrations.Get(ctx, integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !ok {
		ctx.JSON(404, utils.ErrorStr("Integration not found"))
		return
	}

	// Check the integration is public or the user created it
	canActivate, err := dbclient.Client.CustomIntegrationGuilds.CanActivate(ctx, integrationId, userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if !canActivate {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to activate this integration"))
		return
	}

	// Check the secret values are valid
	secrets, err := dbclient.Client.CustomIntegrationSecrets.GetByIntegration(ctx, integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if len(secrets) != len(data.Secrets) {
		ctx.JSON(400, utils.ErrorStr("Invalid secret values"))
		return
	}

	// Since we've checked the length, we can just iterate over the secrets, and they're guaranteed to be correct
	secretMap := make(map[int]string)
	secretValues := make(map[string]string)
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
				secretValues[secret.Name] = value
				break inner
			}
		}

		if !found {
			ctx.JSON(400, utils.ErrorStr("Invalid secret values"))
			return
		}
	}

	// Validate secrets
	if integration.Public && integration.Approved && integration.ValidationUrl != nil {
		integrationHeaders, err := dbclient.Client.CustomIntegrationHeaders.GetByIntegration(ctx, integrationId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		headers := make(map[string]string)
		for _, header := range integrationHeaders {
			value := header.Value
			for key, secret := range secretValues {
				value = strings.ReplaceAll(value, fmt.Sprintf("%%%s%%", key), secret)
			}

			headers[header.Name] = value
		}

		res, statusCode, err := utils.SecureProxyClient.DoRequest(http.MethodPost, *integration.ValidationUrl, headers, secretValues)
		if err != nil {
			if statusCode == http.StatusRequestTimeout {
				ctx.JSON(400, utils.ErrorStr("Secret validation server did not respond in time (contact the integration author)"))
				return
			} else {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		}

		type validationResponse struct {
			Error string `json:"error"`
		}

		var useClientError bool
		var parsed validationResponse
		if err := json.Unmarshal(res, &parsed); err == nil {
			useClientError = len(parsed.Error) > 0

			if len(parsed.Error) > 255 {
				parsed.Error = parsed.Error[:255]
			}
		}

		if statusCode > 299 {
			if useClientError {
				ctx.JSON(400, gin.H{
					"success":      false,
					"error":        "Integration rejected the secret values (contact the integration author for help)",
					"client_error": parsed.Error,
				})
			} else {
				ctx.JSON(400, utils.ErrorStr("Integration rejected the secret values (contact the integration author for help)"))
			}

			return
		}
	}

	if err := dbclient.Client.CustomIntegrationGuilds.AddToGuildWithSecrets(ctx, integrationId, guildId, secretMap); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

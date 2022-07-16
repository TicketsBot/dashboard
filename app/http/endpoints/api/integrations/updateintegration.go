package api

import (
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

type integrationUpdateBody struct {
	Name             string  `json:"name" validate:"required,min=1,max=32"`
	Description      string  `json:"description" validate:"required,min=1,max=255"`
	ImageUrl         *string `json:"image_url" validate:"omitempty,url,max=255,startswith=https://"`
	PrivacyPolicyUrl *string `json:"privacy_policy_url" validate:"omitempty,url,max=255,startswith=https://"`

	Method     string `json:"http_method" validate:"required,oneof=GET POST"`
	WebhookUrl string `json:"webhook_url" validate:"required,webhook,max=255"`

	Secrets []struct {
		Id   int    `json:"id" validate:"omitempty,min=1"`
		Name string `json:"name" validate:"required,min=1,max=32,excludesall=% "`
	} `json:"secrets" validate:"dive,omitempty,min=0,max=5"`

	Headers []struct {
		Id    int    `json:"id" validate:"omitempty,min=1"`
		Name  string `json:"name" validate:"required,min=1,max=32,excludes= "`
		Value string `json:"value" validate:"required,min=1,max=255"`
	} `json:"headers" validate:"dive,omitempty,min=0,max=5"`

	Placeholders []struct {
		Id          int    `json:"id" validate:"omitempty,min=1"`
		Placeholder string `json:"name" validate:"required,min=1,max=32,excludesall=% "`
		JsonPath    string `json:"json_path" validate:"required,min=1,max=255"`
	} `json:"placeholders" validate:"dive,omitempty,min=0,max=15"`
}

func UpdateIntegrationHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	integrationId, err := strconv.Atoi(ctx.Param("integrationid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid integration ID"))
		return
	}

	var data integrationUpdateBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if err := validate.Struct(data); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(500, utils.ErrorStr("An error occurred while validating the integration"))
			return
		}

		formatted := "Your input contained the following errors:"
		for _, validationError := range validationErrors {
			formatted += fmt.Sprintf("\n%s", validationError.Error())
		}

		formatted = strings.TrimSuffix(formatted, "\n")
		ctx.JSON(400, utils.ErrorStr(formatted))
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

	if integration.OwnerId != userId {
		ctx.JSON(403, utils.ErrorStr("You do not own this integration"))
		return
	}

	// Update integration metadata
	err = dbclient.Client.CustomIntegrations.Update(database.CustomIntegration{
		Id:               integration.Id,
		OwnerId:          integration.OwnerId,
		HttpMethod:       data.Method,
		WebhookUrl:       data.WebhookUrl,
		Name:             data.Name,
		Description:      data.Description,
		ImageUrl:         data.ImageUrl,
		PrivacyPolicyUrl: data.PrivacyPolicyUrl,
		Public:           integration.Public,
		Approved:         integration.Approved,
	})

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Store secrets
	if !data.updateSecrets(ctx, integration.Id) {
		return
	}

	// Store headers
	if !data.updateHeaders(ctx, integration.Id) {
		return
	}

	// Store placeholders
	if !data.updatePlaceholders(ctx, integration.Id) {
		return
	}

	ctx.JSON(200, integration)
}

func (b *integrationUpdateBody) updatePlaceholders(ctx *gin.Context, integrationId int) bool {
	// Verify IDs are valid for the integration
	existingPlaceholders, err := dbclient.Client.CustomIntegrationPlaceholders.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	for _, placeholder := range b.Placeholders {
		if placeholder.Id != 0 {
			isValid := false
		inner:
			for _, existingPlaceholder := range existingPlaceholders {
				if existingPlaceholder.Id == placeholder.Id {
					if existingPlaceholder.IntegrationId == integrationId {
						isValid = true
						break inner
					} else {
						ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for placeholders"))
						return false
					}
				}
			}

			if !isValid {
				ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for placeholders"))
				return false
			}
		}
	}

	placeholders := make([]database.CustomIntegrationPlaceholder, len(b.Placeholders))
	for i, placeholder := range b.Placeholders {
		placeholders[i] = database.CustomIntegrationPlaceholder{
			Id:       placeholder.Id,
			Name:     placeholder.Placeholder,
			JsonPath: placeholder.JsonPath,
		}
	}

	if _, err := dbclient.Client.CustomIntegrationPlaceholders.Set(integrationId, placeholders); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	return true
}

func (b *integrationUpdateBody) updateHeaders(ctx *gin.Context, integrationId int) bool {
	// Verify IDs are valid for the integration
	existingHeaders, err := dbclient.Client.CustomIntegrationHeaders.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	for _, header := range b.Headers {
		if header.Id != 0 {
			isValid := false

		inner:
			for _, existingHeader := range existingHeaders {
				if existingHeader.Id == header.Id {
					if existingHeader.IntegrationId == integrationId {
						isValid = true
						break inner
					} else {
						ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for headers"))
						return false
					}
				}
			}

			if !isValid {
				ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for headers"))
				return false
			}
		}
	}

	headers := make([]database.CustomIntegrationHeader, len(b.Headers))
	for i, header := range b.Headers {
		headers[i] = database.CustomIntegrationHeader{
			Id:    header.Id,
			Name:  header.Name,
			Value: header.Value,
		}
	}

	if _, err := dbclient.Client.CustomIntegrationHeaders.CreateOrUpdate(integrationId, headers); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	return true
}

func (b *integrationUpdateBody) updateSecrets(ctx *gin.Context, integrationId int) bool {
	// Verify IDs are valid for the integration
	existingSecrets, err := dbclient.Client.CustomIntegrationSecrets.GetByIntegration(integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	for _, secret := range b.Secrets {
		if secret.Id != 0 {
			isValid := false
		inner:
			for _, existingSecret := range existingSecrets {
				if existingSecret.Id == secret.Id {
					if existingSecret.IntegrationId == integrationId {
						isValid = true
						break inner
					} else {
						ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for secrets"))
						return false
					}
				}
			}

			if !isValid {
				ctx.JSON(400, utils.ErrorStr("Integration ID mismatch for secrets"))
				return false
			}
		}
	}

	secrets := make([]database.CustomIntegrationSecret, len(b.Secrets))
	for i, secret := range b.Secrets {
		secrets[i] = database.CustomIntegrationSecret{
			Id:   secret.Id,
			Name: secret.Name,
		}
	}

	if _, err := dbclient.Client.CustomIntegrationSecrets.CreateOrUpdate(integrationId, secrets); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return false
	}

	return true
}

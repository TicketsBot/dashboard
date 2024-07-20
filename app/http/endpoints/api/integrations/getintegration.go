package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type integrationResponse struct {
	// Strip out the sensitive fields
	Id               int     `json:"id"`
	OwnerId          uint64  `json:"owner_id"`
	WebhookHost      string  `json:"webhook_url"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ImageUrl         *string `json:"image_url"`
	ProxyToken       *string `json:"proxy_token,omitempty"`
	PrivacyPolicyUrl *string `json:"privacy_policy_url"`
	Public           bool    `json:"public"`
	Approved         bool    `json:"approved"`

	Placeholders []database.CustomIntegrationPlaceholder `json:"placeholders"`
	Secrets      []database.CustomIntegrationSecret      `json:"secrets"`
}

func GetIntegrationHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	integrationId, err := strconv.Atoi(ctx.Param("integrationid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid integration ID"))
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

	// Check if the user has permission to view this integration
	if integration.OwnerId != userId && !(integration.Public && integration.Approved) {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to view this integration"))
		return
	}

	placeholders, err := dbclient.Client.CustomIntegrationPlaceholders.GetByIntegration(ctx, integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Don't serve null
	if placeholders == nil {
		placeholders = make([]database.CustomIntegrationPlaceholder, 0)
	}

	secrets, err := dbclient.Client.CustomIntegrationSecrets.GetByIntegration(ctx, integrationId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Don't serve null
	if secrets == nil {
		secrets = make([]database.CustomIntegrationSecret, 0)
	}

	var proxyToken *string
	if integration.ImageUrl != nil {
		tmp, err := utils.GenerateImageProxyToken(*integration.ImageUrl)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		proxyToken = &tmp
	}

	ctx.JSON(200, integrationResponse{
		Id:               integration.Id,
		OwnerId:          integration.OwnerId,
		WebhookHost:      utils.SecondLevelDomain(utils.GetUrlHost(strings.ReplaceAll(integration.WebhookUrl, "%", ""))),
		Name:             integration.Name,
		Description:      integration.Description,
		ImageUrl:         integration.ImageUrl,
		ProxyToken:       proxyToken,
		PrivacyPolicyUrl: integration.PrivacyPolicyUrl,
		Public:           integration.Public,
		Approved:         integration.Approved,
		Placeholders:     placeholders,
		Secrets:          secrets,
	})
}

package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

const pageLimit = 20
const builtInCount = 1

type integrationWithMetadata struct {
	integrationResponse
	GuildCount int  `json:"guild_count"`
	Added      bool `json:"added"`
}

func ListIntegrationsHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)
	guildId := ctx.Keys["guildid"].(uint64)

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 1 {
		page = 1
	}

	page -= 1

	limit := pageLimit
	if page == 0 {
		limit -= builtInCount
	}

	availableIntegrations, err := dbclient.Client.CustomIntegrationGuilds.GetAvailableIntegrationsWithActive(guildId, userId, limit, page*pageLimit)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	integrations := make([]integrationWithMetadata, len(availableIntegrations))
	for i, integration := range availableIntegrations {
		var proxyToken *string
		if integration.ImageUrl != nil {
			tmp, err := utils.GenerateImageProxyToken(*integration.ImageUrl)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			proxyToken = &tmp
		}

		integrations[i] = integrationWithMetadata{
			integrationResponse: integrationResponse{
				Id:               integration.Id,
				OwnerId:          integration.OwnerId,
				WebhookHost:      utils.GetUrlHost(integration.WebhookUrl),
				Name:             integration.Name,
				Description:      integration.Description,
				ImageUrl:         integration.ImageUrl,
				ProxyToken:       proxyToken,
				PrivacyPolicyUrl: integration.PrivacyPolicyUrl,
				Public:           integration.Public,
				Approved:         integration.Approved,
			},
			GuildCount: integration.GuildCount,
			Added:      integration.Active,
		}
	}

	// Don't serve null
	if integrations == nil {
		integrations = make([]integrationWithMetadata, 0)
	}

	ctx.JSON(200, integrations)
}

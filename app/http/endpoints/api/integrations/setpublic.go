package api

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/rest"
	"strconv"
)

func SetIntegrationPublicHandler(ctx *gin.Context) {
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

	if integration.OwnerId != userId {
		ctx.JSON(403, utils.ErrorStr("You do not have permission to manage this integration"))
		return
	}

	if integration.Public {
		ctx.JSON(400, utils.ErrorStr("You have already requested to make this integration public"))
		return
	}

	e := embed.NewEmbed().
		SetTitle("Public Integration Request").
		SetColor(0xfcb97d).
		AddField("Integration ID", strconv.Itoa(integration.Id), true).
		AddField("Integration Name", integration.Name, true).
		AddField("Integration URL", fmt.Sprintf("`%s`", integration.WebhookUrl), true).
		AddField("Integration Owner", fmt.Sprintf("<@%d>", integration.OwnerId), true).
		AddField("Integration Description", integration.Description, false)

	botCtx := botcontext.PublicContext()

	// TODO: Use proper context
	_, err = rest.ExecuteWebhook(
		context.Background(),
		config.Conf.Bot.PublicIntegrationRequestWebhookToken,
		botCtx.RateLimiter,
		config.Conf.Bot.PublicIntegrationRequestWebhookId,
		true,
		rest.WebhookBody{
			Embeds: utils.Slice(e),
		},
	)

	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if err := dbclient.Client.CustomIntegrations.SetPublic(integration.Id); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

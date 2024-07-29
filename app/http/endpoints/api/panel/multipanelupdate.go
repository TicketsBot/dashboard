package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func MultiPanelUpdate(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	// parse body
	var data multiPanelCreateData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorJson(err))
		return
	}

	// parse panel ID
	panelId, err := strconv.Atoi(c.Param("panelid"))
	if err != nil {
		c.JSON(400, utils.ErrorJson(err))
		return
	}

	// retrieve panel from DB
	multiPanel, ok, err := dbclient.Client.MultiPanels.Get(c, panelId)
	if err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// check panel exists
	if !ok {
		c.JSON(404, utils.ErrorJson(errors.New("No panel with the provided ID found")))
		return
	}

	// check panel is in the same guild
	if guildId != multiPanel.GuildId {
		c.JSON(403, utils.ErrorJson(errors.New("Guild ID doesn't match")))
		return
	}

	if err := validate.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); !ok {
			c.JSON(500, utils.ErrorStr("An error occurred while validating the panel"))
			return
		}

		formatted := "Your input contained the following errors:\n" + utils.FormatValidationErrors(validationErrors)
		c.JSON(400, utils.ErrorStr(formatted))
		return
	}

	// validate body & get sub-panels
	panels, err := data.doValidations(guildId)
	if err != nil {
		c.JSON(400, utils.ErrorJson(err))
		return
	}

	for _, panel := range panels {
		if panel.CustomId == "" {
			panel.CustomId, err = utils.RandString(30)
			if err != nil {
				c.JSON(500, utils.ErrorJson(err))
				return
			}

			if err := dbclient.Client.Panel.Update(c, panel); err != nil {
				c.JSON(500, utils.ErrorJson(err))
				return
			}
		}
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// delete old message
	ctx, cancel := app.DefaultContext()
	defer cancel()

	var unwrapped request.RestError
	if err := rest.DeleteMessage(ctx, botContext.Token, botContext.RateLimiter, multiPanel.ChannelId, multiPanel.MessageId); err != nil && !(errors.As(err, &unwrapped) && unwrapped.IsClientError()) {
		c.JSON(500, utils.ErrorJson(err))
		return
	}
	cancel()

	// get premium status
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// send new message
	messageData := data.IntoMessageData(premiumTier > premium.None)
	messageId, err := messageData.send(botContext, panels)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			c.JSON(500, utils.ErrorJson(errors.New("I do not have permission to send messages in the provided channel")))
		} else {
			c.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	// update DB
	dbEmbed, dbEmbedFields := data.Embed.IntoDatabaseStruct()
	updated := database.MultiPanel{
		Id:                    multiPanel.Id,
		MessageId:             messageId,
		ChannelId:             data.ChannelId,
		GuildId:               guildId,
		SelectMenu:            data.SelectMenu,
		SelectMenuPlaceholder: data.SelectMenuPlaceholder,
		Embed: &database.CustomEmbedWithFields{
			CustomEmbed: dbEmbed,
			Fields:      dbEmbedFields,
		},
	}

	if err = dbclient.Client.MultiPanels.Update(c, multiPanel.Id, updated); err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// TODO: one query for ACID purposes
	// delete old targets
	if err := dbclient.Client.MultiPanelTargets.DeleteAll(c, multiPanel.Id); err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	// insert new targets
	group, _ := errgroup.WithContext(context.Background())
	for _, panel := range panels {
		panel := panel

		group.Go(func() error {
			return dbclient.Client.MultiPanelTargets.Insert(c, multiPanel.Id, panel.PanelId)
		})
	}

	if err := group.Wait(); err != nil {
		c.JSON(500, utils.ErrorJson(err))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    multiPanel,
	})
}

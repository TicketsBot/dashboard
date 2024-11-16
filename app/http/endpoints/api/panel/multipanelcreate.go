package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type multiPanelCreateData struct {
	ChannelId             uint64             `json:"channel_id,string"`
	SelectMenu            bool               `json:"select_menu"`
	SelectMenuPlaceholder *string            `json:"select_menu_placeholder,omitempty" validate:"omitempty,max=150"`
	Panels                []int              `json:"panels"`
	Embed                 *types.CustomEmbed `json:"embed" validate:"omitempty,dive"`
}

func (d *multiPanelCreateData) IntoMessageData(isPremium bool) multiPanelMessageData {
	return multiPanelMessageData{
		IsPremium:             isPremium,
		ChannelId:             d.ChannelId,
		SelectMenu:            d.SelectMenu,
		SelectMenuPlaceholder: d.SelectMenuPlaceholder,
		Embed:                 d.Embed.IntoDiscordEmbed(),
	}
}

func MultiPanelCreate(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	var data multiPanelCreateData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorJson(err))
		return
	}

	if err := validate.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewError(err, "An error occurred while validating the panel"))
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

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// get premium status
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(c, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	messageData := data.IntoMessageData(premiumTier > premium.None)
	messageId, err := messageData.send(botContext, panels)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped); unwrapped.StatusCode == 403 {
			c.JSON(http.StatusBadRequest, utils.ErrorJson(errors.New("I do not have permission to send messages in the provided channel")))
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		}

		return
	}

	dbEmbed, dbEmbedFields := data.Embed.IntoDatabaseStruct()
	multiPanel := database.MultiPanel{
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

	multiPanel.Id, err = dbclient.Client.MultiPanels.Create(c, multiPanel)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	group, _ := errgroup.WithContext(context.Background())
	for _, panel := range panels {
		panel := panel

		group.Go(func() error {
			return dbclient.Client.MultiPanelTargets.Insert(c, multiPanel.Id, panel.PanelId)
		})
	}

	if err := group.Wait(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    multiPanel,
	})
}

func (d *multiPanelCreateData) doValidations(guildId uint64) (panels []database.Panel, err error) {
	if err := validateEmbed(d.Embed); err != nil {
		return nil, err
	}

	group, _ := errgroup.WithContext(context.Background())

	group.Go(d.validateChannel(guildId))
	group.Go(func() (e error) {
		panels, e = d.validatePanels(guildId)
		return
	})

	err = group.Wait()
	return
}

func (d *multiPanelCreateData) validateChannel(guildId uint64) func() error {
	return func() error {
		// TODO: Use proper context
		channels, err := cache.Instance.GetGuildChannels(context.Background(), guildId)
		if err != nil {
			return err
		}

		var valid bool
		for _, ch := range channels {
			if ch.Id == d.ChannelId && (ch.Type == channel.ChannelTypeGuildText || ch.Type == channel.ChannelTypeGuildNews) {
				valid = true
				break
			}
		}

		if !valid {
			return errors.New("channel does not exist")
		}

		return nil
	}
}

func (d *multiPanelCreateData) validatePanels(guildId uint64) (panels []database.Panel, err error) {
	if len(d.Panels) < 2 {
		err = errors.New("a multi-panel must contain at least 2 sub-panels")
		return
	}

	if len(d.Panels) > 15 {
		err = errors.New("multi-panels cannot contain more than 15 sub-panels")
		return
	}

	existingPanels, err := dbclient.Client.Panel.GetByGuild(context.Background(), guildId)
	if err != nil {
		return nil, err
	}

	for _, panelId := range d.Panels {
		var valid bool
		// find panel struct
		for _, panel := range existingPanels {
			if panel.PanelId == panelId {
				valid = true
				panels = append(panels, panel)
			}
		}

		if !valid {
			return nil, errors.New("invalid panel ID")
		}
	}

	return
}

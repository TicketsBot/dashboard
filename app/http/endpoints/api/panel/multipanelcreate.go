package api

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
)

type multiPanelCreateData struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Colour     int32  `json:"colour"`
	ChannelId  uint64 `json:"channel_id,string"`
	SelectMenu bool   `json:"select_menu"`
	Panels     []int  `json:"panels"`
}

func (d *multiPanelCreateData) IntoMessageData(isPremium bool) multiPanelMessageData {
	return multiPanelMessageData{
		ChannelId:  d.ChannelId,
		Title:      d.Title,
		Content:    d.Content,
		Colour:     int(d.Colour),
		SelectMenu: d.SelectMenu,
		IsPremium:  isPremium,
	}
}

func MultiPanelCreate(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var data multiPanelCreateData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// validate body & get sub-panels
	panels, err := data.doValidations(guildId)
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get bot context
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// get premium status
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	messageData := data.IntoMessageData(premiumTier > premium.None)
	messageId, err := messageData.send(&botContext, panels)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped); unwrapped.StatusCode == 403 {
			ctx.JSON(500, utils.ErrorJson(errors.New("I do not have permission to send messages in the provided channel")))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	multiPanel := database.MultiPanel{
		MessageId:  messageId,
		ChannelId:  data.ChannelId,
		GuildId:    guildId,
		Title:      data.Title,
		Content:    data.Content,
		Colour:     int(data.Colour),
		SelectMenu: data.SelectMenu,
	}

	multiPanel.Id, err = dbclient.Client.MultiPanels.Create(multiPanel)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	group, _ := errgroup.WithContext(context.Background())
	for _, panel := range panels {
		panel := panel

		group.Go(func() error {
			return dbclient.Client.MultiPanelTargets.Insert(multiPanel.Id, panel.PanelId)
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    multiPanel,
	})
}

func (d *multiPanelCreateData) doValidations(guildId uint64) (panels []database.Panel, err error) {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(d.validateTitle)
	group.Go(d.validateContent)
	group.Go(d.validateChannel(guildId))
	group.Go(func() (e error) {
		panels, e = d.validatePanels(guildId)
		return
	})

	err = group.Wait()
	return
}

func (d *multiPanelCreateData) validateTitle() (err error) {
	if len(d.Title) > 255 {
		err = errors.New("Embed title must be between 1 and 255 characters")
	} else if len(d.Title) == 0 {
		d.Title = "Click to open a ticket"
	}

	return
}

func (d *multiPanelCreateData) validateContent() (err error) {
	if len(d.Content) > 1024 {
		err = errors.New("Embed content must be between 1 and 1024 characters")
	} else if len(d.Content) == 0 { // Fill default
		d.Content = "Click on the button corresponding to the type of ticket you wish to open"
	}

	return
}

func (d *multiPanelCreateData) validateChannel(guildId uint64) func() error {
	return func() (err error) {
		channels := cache.Instance.GetGuildChannels(guildId)

		var valid bool
		for _, ch := range channels {
			if ch.Id == d.ChannelId && ch.Type == channel.ChannelTypeGuildText {
				valid = true
				break
			}
		}

		if !valid {
			err = errors.New("channel does not exist")
		}

		return
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

	existingPanels, err := dbclient.Client.Panel.GetByGuild(guildId)
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

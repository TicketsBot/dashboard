package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func MultiPanelList(ctx *gin.Context) {
	type multiPanelResponse struct {
		database.MultiPanel
		Panels []database.Panel `json:"panels"`
	}

	guildId := ctx.Keys["guildid"].(uint64)

	multiPanels, err := dbclient.Client.MultiPanels.GetByGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorToResponse(err))
		return
	}

	data := make([]multiPanelResponse, len(multiPanels))
	group, _ := errgroup.WithContext(context.Background())
	for i, multiPanel := range multiPanels {
		i := i
		multiPanel := multiPanel

		data[i] = multiPanelResponse{
			MultiPanel: multiPanel,
		}

		group.Go(func() error {
			panels, err := dbclient.Client.MultiPanelTargets.GetPanels(multiPanel.Id)
			if err != nil {
				return err
			}

			data[i].Panels = panels
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorToResponse(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data": data,
	})
}

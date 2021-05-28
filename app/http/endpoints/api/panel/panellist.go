package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func ListPanels(ctx *gin.Context) {
	type panelResponse struct {
		PanelId         int                    `json:"panel_id"`
		ChannelId       uint64                 `json:"channel_id,string"`
		Title           string                 `json:"title"`
		Content         string                 `json:"content"`
		Colour          uint32                 `json:"colour"`
		CategoryId      uint64                 `json:"category_id,string"`
		Emote           string                 `json:"emote"`
		WelcomeMessage  *string                `json:"welcome_message"`
		Mentions        []string               `json:"mentions"`
		WithDefaultTeam bool                   `json:"default_team"`
		Teams           []database.SupportTeam `json:"teams"`
	}

	guildId := ctx.Keys["guildid"].(uint64)

	panels, err := dbclient.Client.Panel.GetByGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	wrapped := make([]panelResponse, len(panels))

	// we will need to lookup role mentions
	group, _ := errgroup.WithContext(context.Background())

	for i, p := range panels {
		i := i
		p := p

		group.Go(func() error {
			var mentions []string

			// get role mentions
			roles, err := dbclient.Client.PanelRoleMentions.GetRoles(p.PanelId)
			if err != nil {
				return err
			}

			// convert to strings
			for _, roleId := range roles {
				mentions = append(mentions, strconv.FormatUint(roleId, 10))
			}

			// get if we should mention the ticket opener
			shouldMention, err := dbclient.Client.PanelUserMention.ShouldMentionUser(p.PanelId)
			if err != nil {
				return err
			}

			if shouldMention {
				mentions = append(mentions, "user")
			}

			teams, err := dbclient.Client.PanelTeams.GetTeams(p.MessageId)
			if err != nil {
				return err
			}

			wrapped[i] = panelResponse{
				PanelId:         p.PanelId,
				ChannelId:       p.ChannelId,
				Title:           p.Title,
				Content:         p.Content,
				Colour:          uint32(p.Colour),
				CategoryId:      p.TargetCategory,
				Emote:           p.ReactionEmote,
				WelcomeMessage:  p.WelcomeMessage,
				Mentions:        mentions,
				WithDefaultTeam: p.WithDefaultTeam,
				Teams:           teams,
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, wrapped)
}

package api

import (
	"context"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func ListPanels(ctx *gin.Context) {
	type panelResponse struct {
		database.Panel
		UseCustomEmoji               bool        `json:"use_custom_emoji"`
		Emoji                        types.Emoji `json:"emote"`
		Mentions                     []string    `json:"mentions"`
		Teams                        []int       `json:"teams"`
		UseServerDefaultNamingScheme bool        `json:"use_server_default_naming_scheme"`
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

			// get if we should mention the ticket opener
			shouldMention, err := dbclient.Client.PanelUserMention.ShouldMentionUser(p.PanelId)
			if err != nil {
				return err
			}

			if shouldMention {
				mentions = append(mentions, "user")
			}

			// get role mentions
			roles, err := dbclient.Client.PanelRoleMentions.GetRoles(p.PanelId)
			if err != nil {
				return err
			}

			// convert to strings
			for _, roleId := range roles {
				mentions = append(mentions, strconv.FormatUint(roleId, 10))
			}

			teamIds, err := dbclient.Client.PanelTeams.GetTeamIds(p.PanelId)
			if err != nil {
				return err
			}

			// Don't serve null
			if teamIds == nil {
				teamIds = make([]int, 0)
			}

			wrapped[i] = panelResponse{
				Panel:                        p,
				UseCustomEmoji:               p.EmojiId != nil,
				Emoji:                        types.NewEmoji(p.EmojiName, p.EmojiId),
				Mentions:                     mentions,
				Teams:                        teamIds,
				UseServerDefaultNamingScheme: p.NamingScheme == nil,
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

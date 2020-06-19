package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"strconv"
)

type panel struct {
	ChannelId      uint64   `json:"channel_id,string"`
	MessageId      uint64   `json:"message_id,string"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Colour         uint32   `json:"colour"`
	CategoryId     uint64   `json:"category_id,string"`
	Emote          string   `json:"emote"`
	WelcomeMessage *string  `json:"welcome_message"`
	Mentions       []string `json:"mentions"`
}

func ListPanels(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	panels, err := database.Client.Panel.GetByGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	wrapped := make([]panel, len(panels))

	// we will need to lookup role mentions
	group, _ := errgroup.WithContext(context.Background())

	for i, p := range panels {
		group.Go(func() (err error) {
			var mentions []string

			// get role mentions
			roles, err := database.Client.PanelRoleMentions.GetRoles(p.MessageId)
			if err != nil {
				return err
			}

			// convert to strings
			for _, roleId := range roles {
				mentions = append(mentions, strconv.FormatUint(roleId, 10))
			}

			// get if we should mention the ticket opener
			shouldMention, err := database.Client.PanelUserMention.ShouldMentionUser(p.MessageId)
			if err != nil {
				return err
			}

			if shouldMention {
				mentions = append(mentions, "user")
			}

			wrapped[i] = panel{
				MessageId:      p.MessageId,
				ChannelId:      p.ChannelId,
				Title:          p.Title,
				Content:        p.Content,
				Colour:         uint32(p.Colour),
				CategoryId:     p.TargetCategory,
				Emote:          p.ReactionEmote,
				WelcomeMessage: p.WelcomeMessage,
				Mentions:       mentions,
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, wrapped)
}

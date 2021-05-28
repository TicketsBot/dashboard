package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
)

const freePanelLimit = 3

type panelBody struct {
	ChannelId      uint64   `json:"channel_id,string"`
	MessageId      uint64   `json:"message_id,string"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Colour         uint32   `json:"colour"`
	CategoryId     uint64   `json:"category_id,string"`
	Emote          string   `json:"emote"`
	WelcomeMessage *string  `json:"welcome_message"`
	Mentions       []string `json:"mentions"`
	Teams          []string `json:"teams"`
}

func CreatePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var data panelBody

	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	data.MessageId = 0

	// Check panel quota
	premiumTier := rpc.PremiumClient.GetTierByGuildId(guildId, false, botContext.Token, botContext.RateLimiter)

	if premiumTier == premium.None {
		panels, err := dbclient.Client.Panel.GetByGuild(guildId)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		if len(panels) >= freePanelLimit {
			ctx.AbortWithStatusJSON(402, gin.H{
				"success": false,
				"error":   "You have exceeded your panel quota. Purchase premium to unlock more panels.",
			})
			return
		}
	}

	if !data.doValidations(ctx, guildId) {
		return
	}

	customId := utils.RandString(80)

	emoji, _ := data.getEmoji() // already validated
	msgId, err := data.sendEmbed(&botContext, data.Title, customId, emoji, premiumTier > premium.None)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   "I do not have permission to send messages in the specified channel",
			})
		} else {
			// TODO: Most appropriate error?
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		return
	}

	// Store in DB
	panel := database.Panel{
		MessageId:       msgId,
		ChannelId:       data.ChannelId,
		GuildId:         guildId,
		Title:           data.Title,
		Content:         data.Content,
		Colour:          int32(data.Colour),
		TargetCategory:  data.CategoryId,
		ReactionEmote:   emoji,
		WelcomeMessage:  data.WelcomeMessage,
		WithDefaultTeam: utils.ContainsString(data.Teams, "default"),
		CustomId:        customId,
	}

	panelId, err := dbclient.Client.Panel.Create(panel)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// insert role mention data
	// string is role ID or "user" to mention the ticket opener
	for _, mention := range data.Mentions {
		if mention == "user" {
			if err = dbclient.Client.PanelUserMention.Set(panelId, true); err != nil {
				ctx.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			// should we check the role is a valid role in the guild?
			// not too much of an issue if it isnt

			if err = dbclient.Client.PanelRoleMentions.Add(panelId, roleId); err != nil {
				ctx.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}
		}
	}

	if responseCode, err := insertTeams(guildId, msgId, data.Teams); err != nil {
		ctx.JSON(responseCode, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success":  true,
		"panel_id": panelId,
	})
}

// returns (response_code, error)
func insertTeams(guildId, panelMessageId uint64, teamIds []string) (int, error) {
	// insert teams
	group, _ := errgroup.WithContext(context.Background())
	for _, teamId := range teamIds {
		if teamId == "default" {
			continue // already handled
		}

		teamId, err := strconv.Atoi(teamId)
		if err != nil {
			return 400, err
		}

		group.Go(func() error {
			// ensure team exists
			exists, err := dbclient.Client.SupportTeam.Exists(teamId, guildId)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("team with id %d not found", teamId)
			}

			return dbclient.Client.PanelTeams.Add(panelMessageId, teamId)
		})
	}

	return 500, group.Wait()
}

func (p *panelBody) doValidations(ctx *gin.Context, guildId uint64) bool {
	if !p.verifyTitle() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel titles must be between 1 - 80 characters in length",
		})
		return false
	}

	if !p.verifyContent() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel content must be between 1 - 1024 characters in length",
		})
		return false
	}

	channels := cache.Instance.GetGuildChannels(guildId)

	if !p.verifyChannel(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel",
		})
		return false
	}

	if !p.verifyCategory(channels) {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid channel category",
		})
		return false
	}

	_, validEmoji := p.getEmoji()
	if !validEmoji {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Invalid emoji. Simply use the emoji's name from Discord.",
		})
		return false
	}

	if !p.verifyWelcomeMessage() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Welcome message must be null or between 1 - 1024 characters",
		})
		return false
	}

	return true
}

func (p *panelBody) verifyTitle() bool {
	return len(p.Title) > 0 && len(p.Title) <= 80
}

func (p *panelBody) verifyContent() bool {
	return len(p.Content) > 0 && len(p.Content) < 1025
}

func (p *panelBody) getEmoji() (emoji string, ok bool) {
	p.Emote = strings.Replace(p.Emote, ":", "", -1)

	emoji, ok = utils.GetEmoji(p.Emote)
	return
}

func (p *panelBody) verifyChannel(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.ChannelId && ch.Type == channel.ChannelTypeGuildText {
			valid = true
			break
		}
	}

	return valid
}

func (p *panelBody) verifyCategory(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.CategoryId && ch.Type == channel.ChannelTypeGuildCategory {
			valid = true
			break
		}
	}

	return valid
}

func (p *panelBody) verifyWelcomeMessage() bool {
	return p.WelcomeMessage == nil || (len(*p.WelcomeMessage) > 0 && len(*p.WelcomeMessage) < 1025)
}

func (p *panelBody) sendEmbed(ctx *botcontext.BotContext, title, customId, emote string, isPremium bool) (uint64, error) {
	e := embed.NewEmbed().
		SetTitle(p.Title).
		SetDescription(p.Content).
		SetColor(int(p.Colour))

	if !isPremium {
		// TODO: Don't harcode
		e.SetFooter("Powered by ticketsbot.net", "https://cdn.discordapp.com/avatars/508391840525975553/ac2647ffd4025009e2aa852f719a8027.png?size=256")
	}

	data := rest.CreateMessageData{
		Embed: e,
		Components: []component.Component{
			{
				Type: component.ComponentActionRow,
				ComponentData: component.ActionRow{
					Components: []component.Component{
						{
							Type: component.ComponentButton,
							ComponentData: component.Button{
								Label:    title,
								CustomId: customId,
								Style:    component.ButtonStylePrimary,
								Emoji: emoji.Emoji{
									Name: emote,
								},
								Url:      nil,
								Disabled: false,
							},
						},
					},
				},
			},
		},
	}

	msg, err := rest.CreateMessage(ctx.Token, ctx.RateLimiter, p.ChannelId, data)
	if err != nil {
		return 0, err
	}

	return msg.Id, nil
}

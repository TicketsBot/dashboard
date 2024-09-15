package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
	"strconv"
)

type sendMessageBody struct {
	MessageType string `json:"type"`
	Content     string `json:"content"`
}

func SendMessage(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Get ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid ticket ID",
		})
		return
	}

	var body sendMessageBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Message is missing",
		})
		return
	}

	if len(body.Content) == 0 {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "You must enter a message",
		})
		return
	}

	// Verify guild is premium
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if premiumTier == premium.None {
		ctx.JSON(402, gin.H{
			"success": false,
			"error":   "Guild is not premium",
		})
		return
	}

	// Get ticket
	ticket, err := database.Client.Tickets.Get(ctx, ticketId, guildId)

	// Verify the ticket exists
	if ticket.UserId == 0 {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Ticket not found",
		})
		return
	}

	// Verify the user has permission to send to this guild
	if ticket.GuildId != guildId {
		ctx.JSON(403, gin.H{
			"success": false,
			"error":   "Guild ID doesn't match",
		})
		return
	}

	if len(body.Content) > 2000 {
		body.Content = body.Content[0:1999]
	}

	// Preferably send via a webhook
	webhook, err := database.Client.Webhooks.Get(ctx, guildId, ticketId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	settings, err := database.Client.Settings.Get(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorStr("Failed to fetch settings"))
		return
	}

	if webhook.Id != 0 {
		var webhookData rest.WebhookBody
		if settings.AnonymiseDashboardResponses {
			guild, err := botContext.GetGuild(context.Background(), guildId)
			if err != nil {
				ctx.JSON(500, utils.ErrorStr("Failed to fetch guild"))
				return
			}

			webhookData = rest.WebhookBody{
				Content:   body.Content,
				Username:  guild.Name,
				AvatarUrl: guild.IconUrl(),
			}
		} else {
			user, err := botContext.GetUser(context.Background(), userId)
			if err != nil {
				ctx.JSON(500, utils.ErrorStr("Failed to fetch user"))
				return
			}

			webhookData = rest.WebhookBody{
				Content:   body.Content,
				Username:  user.EffectiveName(),
				AvatarUrl: user.AvatarUrl(256),
			}
		}

		// TODO: Ratelimit
		_, err = rest.ExecuteWebhook(ctx, webhook.Token, nil, webhook.Id, true, webhookData)

		if err != nil {
			// We can delete the webhook in this case
			var unwrapped request.RestError
			if errors.As(err, &unwrapped); unwrapped.StatusCode == 403 || unwrapped.StatusCode == 404 {
				go database.Client.Webhooks.Delete(ctx, guildId, ticketId)
			}
		} else {
			ctx.JSON(200, gin.H{
				"success": true,
			})
			return
		}
	}

	message := body.Content
	if !settings.AnonymiseDashboardResponses {
		user, err := botContext.GetUser(context.Background(), userId)
		if err != nil {
			ctx.JSON(500, utils.ErrorStr("Failed to fetch user"))
			return
		}

		message = fmt.Sprintf("**%s**: %s", user.EffectiveName(), message)
	}

	if len(message) > 2000 {
		message = message[0:1999]
	}

	if ticket.ChannelId == nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"error":   "Ticket channel ID is nil",
		})
		return
	}

	if _, err = rest.CreateMessage(ctx, botContext.Token, botContext.RateLimiter, *ticket.ChannelId, rest.CreateMessageData{Content: message}); err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"net/http"
	"time"
)

type sendMessageBody struct {
	Message string `json:"message"`
}

func SendMessage(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	var body sendMessageBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Message is missing",
		})
		return
	}

	// Verify guild is premium
	isPremium := make(chan bool)
	go utils.IsPremiumGuild(guildId, isPremium)
	if !<-isPremium {
		ctx.AbortWithStatusJSON(402, gin.H{
			"success": false,
			"error":   "Guild is not premium",
		})
		return
	}

	// Get ticket
	ticketChan := make(chan table.Ticket)
	go table.GetTicket(ctx.Param("uuid"), ticketChan)
	ticket := <-ticketChan

	// Verify the ticket exists
	if ticket.TicketId == 0 {
		ctx.AbortWithStatusJSON(404, gin.H{
			"success": false,
			"error":   "Ticket not found",
		})
		return
	}

	// Verify the user has permission to send to this guild
	if ticket.Guild != guildId {
		ctx.AbortWithStatusJSON(403, gin.H{
			"success": false,
			"error":   "Guild ID doesn't match",
		})
		return
	}

	user, _ := cache.Instance.GetUser(userId)

	if len(body.Message) > 2000 {
		body.Message = body.Message[0:1999]
	}

	// Preferably send via a webhook
	webhookChan := make(chan *string)
	go table.GetWebhookByUuid(ticket.Uuid, webhookChan)
	webhook := <-webhookChan

	success := false
	if webhook != nil {
		// TODO: Use gdl execute webhook wrapper
		success = executeWebhook(ticket.Uuid, *webhook, body.Message, user.Username, user.AvatarUrl(256))
	}

	if !success {
		body.Message = fmt.Sprintf("**%s**: %s", user.Username, body.Message)
		if len(body.Message) > 2000 {
			body.Message = body.Message[0:1999]
		}

		_, _ = rest.CreateMessage(config.Conf.Bot.Token, ratelimit.Ratelimiter, ticket.Channel, rest.CreateMessageData{Content: body.Message})
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

func executeWebhook(uuid, webhook, content, username, avatar string) bool {
	body := map[string]interface{}{
		"content":    content,
		"username":   username,
		"avatar_url": avatar,
	}
	encoded, err := json.Marshal(&body)
	if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://canary.discordapp.com/api/webhooks/%s", webhook), bytes.NewBuffer(encoded))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	res, err := client.Do(req)
	if err != nil {
		return false
	}

	if res.StatusCode == 404 || res.StatusCode == 403 {
		go table.DeleteWebhookByUuid(uuid)
	} else {
		return true
	}

	return false
}

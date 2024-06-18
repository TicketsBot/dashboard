package api

import (
	"context"
	"encoding/base64"
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/tokenchange"
	"github.com/TicketsBot/common/whitelabeldelete"
	"github.com/TicketsBot/database"
	"github.com/TicketsBot/worker/bot/command/manager"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/application"
	"github.com/rxdn/gdl/rest"
	"strconv"
	"strings"
	"time"
)

func WhitelabelPost() func(*gin.Context) {
	cm := new(manager.CommandManager)
	cm.RegisterCommands()

	return func(ctx *gin.Context) {
		userId := ctx.Keys["userid"].(uint64)

		type whitelabelPostBody struct {
			Token string `json:"token"`
		}

		// Get token
		var data whitelabelPostBody
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"error":   "Missing token",
			})
			return
		}

		bot, err := fetchApplication(data.Token)
		if err != nil {
			ctx.JSON(400, utils.ErrorJson(err))
			return
		}

		// Check if this is a different token
		existing, err := dbclient.Client.Whitelabel.GetByUserId(userId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Take existing whitelabel bot offline, if it is a different bot
		if existing.BotId != 0 && existing.BotId != bot.Id {
			whitelabeldelete.Publish(redis.Client.Client, existing.BotId)
		}

		// Set token in DB so that http-gateway can use it when Discord validates the interactions endpoint
		// TODO: Use a transaction
		if err := dbclient.Client.Whitelabel.Set(database.WhitelabelBot{
			UserId: userId,
			BotId:  bot.Id,
			Token:  data.Token,
		}); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if err := dbclient.Client.WhitelabelKeys.Set(bot.Id, bot.VerifyKey); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		// Allow some time for the database change to be propagated
		time.Sleep(time.Millisecond * 500)

		// Set intents
		var currentFlags application.Flag = 0
		if bot.Flags != nil {
			currentFlags = *bot.Flags
		}

		editData := rest.EditCurrentApplicationData{
			Flags: utils.Ptr(application.BuildFlags(
				currentFlags,
				application.FlagIntentGatewayGuildMembersLimited,
				application.FlagGatewayMessageContentLimited,
			)),
			// TODO: Don't hardcode URL
			InteractionsEndpointUrl: utils.Ptr(fmt.Sprintf("https://gateway.ticketsbot.net/handle/%d", bot.InteractionsEndpointUrl)),
		}

		if _, err := rest.EditCurrentApplication(context.Background(), data.Token, nil, editData); err != nil {
			// TODO: Use a transaction
			if err := dbclient.Client.Whitelabel.Delete(bot.Id); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		tokenChangeData := tokenchange.TokenChangeData{
			Token: data.Token,
			NewId: bot.Id,
			OldId: 0,
		}

		if err := tokenchange.PublishTokenChange(redis.Client.Client, tokenChangeData); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if err := createInteractions(cm, bot.Id, data.Token); err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"bot":     bot,
		})
	}
}

func validateToken(token string) bool {
	split := strings.Split(token, ".")

	// Check for 2 dots
	if len(split) != 3 {
		return false
	}

	// Validate bot ID
	// TODO: We could check the date on the snowflake
	idRaw, err := base64.RawStdEncoding.DecodeString(split[0])
	if err != nil {
		return false
	}

	if _, err := strconv.ParseUint(string(idRaw), 10, 64); err != nil {
		return false
	}

	// Validate time
	if _, err := base64.RawURLEncoding.DecodeString(split[1]); err != nil {
		return false
	}

	return true
}

func fetchApplication(token string) (*application.Application, error) {
	if !validateToken(token) {
		return nil, fmt.Errorf("Invalid token")
	}

	// Validate token + get bot ID
	// TODO: Use proper context
	app, err := rest.GetCurrentApplication(context.Background(), token, nil)
	if err != nil {
		return nil, err
	}

	if app.Id == 0 {
		return nil, fmt.Errorf("Invalid token")
	}

	return &app, nil
}

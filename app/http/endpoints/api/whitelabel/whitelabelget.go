package api

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest"
)

type whitelabelResponse struct {
	Id        uint64 `json:"id,string"`
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
	statusUpdateBody
}

func WhitelabelGet(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Check if this is a different token
	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if bot.BotId == 0 {
		ctx.JSON(404, utils.ErrorStr("No bot found"))
		return
	}

	// Get public key
	publicKey, err := database.Client.WhitelabelKeys.Get(bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Get status
	status, statusType, _, err := database.Client.WhitelabelStatuses.Get(bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	username := getBotUsername(context.Background(), bot.Token)

	ctx.JSON(200, whitelabelResponse{
		Id:        bot.BotId,
		PublicKey: publicKey,
		Username:  username,
		statusUpdateBody: statusUpdateBody{ // Zero values if no status is fine
			Status:     status,
			StatusType: user.ActivityType(statusType),
		},
	})
}

func getBotUsername(ctx context.Context, token string) string {
	user, err := rest.GetCurrentUser(ctx, token, nil)
	if err != nil {
		// TODO: Log error
		return "Unknown User"
	}

	return user.Username
}

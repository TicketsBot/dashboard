package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
)

type whitelabelResponse struct {
	Id uint64 `json:"id,string"`
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

	// Get status
	status, statusType, _, err := database.Client.WhitelabelStatuses.Get(bot.BotId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, whitelabelResponse{
		Id: bot.BotId,
		statusUpdateBody: statusUpdateBody{ // Zero values if no status is fine
			Status:     status,
			StatusType: user.ActivityType(statusType),
		},
	})
}

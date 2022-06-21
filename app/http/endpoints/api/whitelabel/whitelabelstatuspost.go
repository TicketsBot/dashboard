package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/statusupdates"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
)

type statusUpdateBody struct {
	Status     string            `json:"status"`
	StatusType user.ActivityType `json:"status_type,string"`
}

func WhitelabelStatusPost(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Get bot
	bot, err := database.Client.Whitelabel.GetByUserId(userId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Ensure bot exists
	if bot.BotId == 0 {
		ctx.JSON(404, utils.ErrorStr("No bot found"))
		return
	}

	// Parse status
	var data statusUpdateBody
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid request body"))
		return
	}

	// Validate status length
	if len(data.Status) == 0 || len(data.Status) > 255 {
		ctx.JSON(400, utils.ErrorStr("Status must be between 1-255 characters in length"))
		return
	}

	// Validate status type
	validActivities := []user.ActivityType{
		user.ActivityTypePlaying,
		user.ActivityTypeListening,
		user.ActivityTypeWatching,
	}

	if !utils.Contains(validActivities, data.StatusType) {
		ctx.JSON(400, utils.ErrorStr("Invalid status type"))
		return
	}

	// Update in database
	if err := database.Client.WhitelabelStatuses.Set(bot.BotId, data.Status, int16(data.StatusType)); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Send status update to sharder
	go statusupdates.Publish(redis.Client.Client, bot.BotId)

	ctx.JSON(200, utils.SuccessResponse)
}

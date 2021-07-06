package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"time"
)

// time.Duration marshals to nanoseconds, custom impl to marshal to seconds
type autoCloseBody struct {
	Enabled                 bool  `json:"enabled"`
	SinceOpenWithNoResponse int64 `json:"since_open_with_no_response"`
	SinceLastMessage        int64 `json:"since_last_message"`
	OnUserLeave             bool  `json:"on_user_leave"`
}

func GetAutoClose(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	settings, err := dbclient.Client.AutoClose.Get(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(200, convertToAutoCloseBody(settings))
}

func convertToAutoCloseBody(settings database.AutoCloseSettings) (body autoCloseBody) {
	body.Enabled = settings.Enabled

	if settings.SinceOpenWithNoResponse != nil {
		body.SinceOpenWithNoResponse = int64(*settings.SinceOpenWithNoResponse / time.Second)
	}

	if settings.SinceLastMessage != nil {
		body.SinceLastMessage = int64(*settings.SinceLastMessage / time.Second)
	}

	if settings.OnUserLeave != nil {
		body.OnUserLeave = *settings.OnUserLeave
	}

	return
}

func convertFromAutoCloseBody(body autoCloseBody) (settings database.AutoCloseSettings) {
	settings.Enabled = body.Enabled

	if body.SinceOpenWithNoResponse > 0 {
		duration := time.Second * time.Duration(body.SinceOpenWithNoResponse)
		settings.SinceOpenWithNoResponse = &duration
	}

	if body.SinceLastMessage > 0 {
		duration := time.Second * time.Duration(body.SinceLastMessage)
		settings.SinceLastMessage = &duration
	}

	settings.OnUserLeave = &body.OnUserLeave

	return
}

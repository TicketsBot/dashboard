package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"time"
)

var maxDays = 90
var maxLength = time.Hour * 24 * time.Duration(maxDays)

func PostAutoClose(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var body autoCloseBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	settings := convertFromAutoCloseBody(body)

	// get premium
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if premiumTier < premium.Premium {
		settings.SinceOpenWithNoResponse = nil
		settings.SinceLastMessage = nil
	}

	// Time period cannot be negative, convertFromAutoCloseBody will not allow
	if (settings.SinceOpenWithNoResponse != nil && *settings.SinceOpenWithNoResponse > maxLength) ||
		(settings.SinceLastMessage != nil && *settings.SinceLastMessage > maxLength) {
		ctx.JSON(400, utils.ErrorStr("Time period cannot be longer than %d days", maxDays))
		return
	}

	if !settings.Enabled {
		settings.SinceLastMessage = nil
		settings.SinceOpenWithNoResponse = nil
		settings.OnUserLeave = nil
	}

	if err := dbclient.Client.AutoClose.Set(guildId, settings); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}

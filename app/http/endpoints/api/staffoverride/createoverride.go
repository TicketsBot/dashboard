package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type createOverrideBody struct {
	TimePeriod int `json:"time_period"`
}

func CreateOverrideHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	var body createOverrideBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, utils.ErrorStr("Invalid request body"))
		fmt.Println(err.Error())
		return
	}

	expires := time.Now().Add(time.Hour * time.Duration(body.TimePeriod))
	if err := database.Client.StaffOverride.Set(guildId, expires); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(204)
}

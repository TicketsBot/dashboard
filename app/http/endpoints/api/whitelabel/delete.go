package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WhitelabelDelete(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	// Check if this is a different token
	if err := database.Client.Whitelabel.Delete(userId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

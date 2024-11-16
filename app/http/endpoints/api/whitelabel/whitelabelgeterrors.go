package api

import (
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WhitelabelGetErrors(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	errors, err := database.Client.WhitelabelErrors.GetRecent(c, userId, 10)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"errors":  errors,
	})
}

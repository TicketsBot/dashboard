package api

import (
	"github.com/TicketsBot/GoPanel/app"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPermissionLevel(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	guildId, err := strconv.ParseUint(c.Query("guild"), 10, 64)
	if err != nil {
		c.JSON(400, utils.ErrorStr("Invalid guild ID"))
		return
	}

	// TODO: Use proper context
	permissionLevel, err := utils.GetPermissionLevel(c, guildId, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, gin.H{
		"success":          true,
		"permission_level": permissionLevel,
	})
}

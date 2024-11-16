package forms

import (
	"github.com/TicketsBot/GoPanel/app"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateForm(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	var data createFormBody
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorJson(err))
		return
	}

	if len(data.Title) > 45 {
		c.JSON(400, utils.ErrorStr("Title is too long"))
		return
	}

	formId, err := strconv.Atoi(c.Param("form_id"))
	if err != nil {
		c.JSON(400, utils.ErrorStr("Invalid form ID"))
		return
	}

	form, ok, err := dbclient.Client.Forms.Get(c, formId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	if !ok {
		c.JSON(404, utils.ErrorStr("Form not found"))
		return
	}

	if form.GuildId != guildId {
		c.JSON(403, utils.ErrorStr("Form does not belong to this guild"))
		return
	}

	if err := dbclient.Client.Forms.UpdateTitle(c, formId, data.Title); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, utils.SuccessResponse)
}

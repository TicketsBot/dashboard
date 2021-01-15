package api

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func TokenHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := utils.GetUserId(session)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": strconv.FormatUint(userId, 10),
		"timestamp": time.Now(),
	})

	str, err := token.SignedString([]byte(config.Conf.Server.Secret))
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"success": true,
			"token": str,
		})
	}
}

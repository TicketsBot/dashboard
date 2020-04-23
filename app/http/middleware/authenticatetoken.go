package middleware

import (
	"errors"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AuthenticateToken(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")

	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Conf.Server.Secret), nil
	})

	if err != nil {
		ctx.AbortWithStatusJSON(403, gin.H{
			"error": err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, hasUserId := claims["userid"]
		if !hasUserId {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": errors.New("token is invalid"),
			})
			return
		}

		parsedId, err := strconv.ParseUint(userId.(string), 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(403, gin.H{
				"error": errors.New("token is invalid"),
			})
			return
		}

		ctx.Keys["userid"] = parsedId
	} else {
		ctx.AbortWithStatusJSON(403, gin.H{
			"error": errors.New("token is invalid"),
		})
		return
	}
}

package utils

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GenerateImageProxyToken(imageUrl string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"url":        imageUrl,
		"request_id": uuid.New().String(),
		"exp":        strconv.FormatInt(time.Now().Add(time.Second*30).Unix(), 10),
	})

	return token.SignedString([]byte(config.Conf.Bot.ImageProxySecret))
}

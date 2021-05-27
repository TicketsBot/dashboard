package utils

import (
	"encoding/base64"
	"math/rand"
	"strconv"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func IsInt(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

func Base64Decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s); if err != nil {
		return ""
	}
	return string(b)
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

package utils

import (
	"github.com/gin-gonic/contrib/sessions"
	"strconv"
)

func IsLoggedIn(store sessions.Session) bool {
	return store.Get("access_token") != nil && store.Get("expiry") != nil && store.Get("refresh_token") != nil && store.Get("userid") != nil && store.Get("name") != nil
}

func GetUserId(store sessions.Session) (int64, error) {
	return strconv.ParseInt(store.Get("userid").(string), 10, 64)
}

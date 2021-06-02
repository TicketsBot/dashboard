package utils

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
)

func IsLoggedIn(store sessions.Session) bool {
	requiredKeys := []string{"access_token", "expiry", "refresh_token", "userid", "name", "avatar", "csrf"}
	for _, key := range requiredKeys {
		if store.Get(key) == nil {
			fmt.Println(key)
			return false
		}
	}

	return true
}

func GetUserId(store sessions.Session) uint64 {
	return store.Get("userid").(uint64)
}

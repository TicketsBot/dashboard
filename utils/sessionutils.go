package utils

import (
	"github.com/gin-gonic/contrib/sessions"
)

func IsLoggedIn(store sessions.Session) bool {
	requiredKeys := []string{"access_token", "expiry", "refresh_token", "userid", "name", "avatar", "csrf"}
	for _, key := range requiredKeys {
		if store.Get(key) == nil {
			return false
		}
	}

	return true
}

func GetUserId(store sessions.Session) uint64 {
	return store.Get("userid").(uint64)
}

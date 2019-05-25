package utils

import (
	"github.com/gin-gonic/contrib/sessions"
)

func IsLoggedIn(store sessions.Session) bool {
	return store.Get("access_token") != nil && store.Get("expiry") != nil && store.Get("refresh_token") != nil && store.Get("userid") != nil && store.Get("name") != nil
}

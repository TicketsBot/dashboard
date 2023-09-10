package utils

import (
	"github.com/TicketsBot/GoPanel/config"
	"net/http"
)

// Twilight's HTTP proxy doesn't support the typical HTTP proxy protocol - instead you send the request directly
// to the proxy's host in the URL. This is not how Go's proxy function should be used, but it works :)
func ProxyHook(token string, req *http.Request) {
	req.URL.Scheme = "http"
	req.URL.Host = config.Conf.Bot.ProxyUrl
}

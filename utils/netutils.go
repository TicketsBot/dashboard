package utils

import (
	"net/url"
)

func GetUrlHost(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "Invalid URL"
	}

	return parsed.Hostname()
}

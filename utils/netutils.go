package utils

import (
	"net/url"
	"strings"
)

func GetUrlHost(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "Invalid URL"
	}

	return parsed.Hostname()
}

func SecondLevelDomain(domain string) string {
	split := strings.Split(strings.TrimRight(domain, "."), ".")
	if len(split) > 2 {
		return strings.Join(split[len(split)-2:len(split)], ".")
	} else {
		return strings.Join(split, ".")
	}
}

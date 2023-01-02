package utils

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"net/url"
)

func GetUrlHost(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "Invalid URL"
	}

	return parsed.Hostname()
}

func SecondLevelDomain(domain string) string {
	domain, err := publicsuffix.Domain(domain)
	if err != nil {
		return "Invalid domain"
	}

	return domain
}

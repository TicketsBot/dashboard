package api

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/go-playground/validator/v10"
	"net/url"
	"regexp"
)

var placeholderRegex = regexp.MustCompile(`%[\w|-]+%`)

func newIntegrationValidator() *validator.Validate {
	v := validator.New()
	utils.Must(v.RegisterValidation("webhook", WebhookValidator))
	return v
}

func WebhookValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	stripped := placeholderRegex.ReplaceAllString(value, "")

	parsed, err := url.Parse(stripped)
	if err != nil {
		return false
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	if parsed.Host == "" {
		return false
	}

	return true
}

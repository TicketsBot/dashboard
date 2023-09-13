package livechat

import (
	"encoding/json"
)

type (
	EventType string

	Event struct {
		Type EventType       `json:"type"`
		Data json.RawMessage `json:"data,omitempty"`
	}

	AuthData struct {
		Token string `json:"token"`
	}

	ErrorMessage struct {
		Error string `json:"error"`
	}
)

const (
	EventTypeAuth          EventType = "auth"
	EventTypeAuthenticated EventType = "authenticated"
	EventTypeMessage       EventType = "message"
)

func NewErrorMessage(message string) ErrorMessage {
	return ErrorMessage{message}
}

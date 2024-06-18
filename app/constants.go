package app

import (
	"context"
	"time"
)

// DefaultTimeout is the default timeout for operations requiring REST requests: this is a high value, as the HTTP proxy
// may delay requests to avoid rate limits.
const DefaultTimeout = time.Second * 30

func DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

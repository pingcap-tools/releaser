package utils

import (
	"context"
	"time"
)

const timeout = 5 * time.Second

// NewTimeoutContext create context with timeout
func NewTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

package middleware

import (
	"context"
)

type contextKey string

var (
	userIdContextKey    contextKey = "userId"
	requestIdContextKey contextKey = "requestId"
)

func withUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, userIdContextKey, userID)
}

func withRequestID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, userIdContextKey, requestID)
}

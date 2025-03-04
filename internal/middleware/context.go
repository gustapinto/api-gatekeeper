package middleware

import "context"

type contextKey string

var (
	userIdContextKey contextKey = "userId"
)

func withUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, userIdContextKey, userID)
}

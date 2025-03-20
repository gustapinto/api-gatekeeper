package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/model"
)

type GuardBackendRouteNextFunc = func(http.ResponseWriter, *http.Request, config.Backend, config.Route)

type GuardApplicationRouteNextFunc = http.HandlerFunc

type AuthService interface {
	AuthenticateToken(string) (model.User, error)

	Authorize(model.User, []string) error
}

type contextKey string

var (
	userIdContextKey    contextKey = "userId"
	requestIdContextKey contextKey = "requestId"
)

func withUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, userIdContextKey, userID)
}

func withRequestID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, requestIdContextKey, requestID)
}

func mergeScopes(backend config.Backend, route config.Route) []string {
	scopes := make([]string, 0)
	scopes = append(scopes, backend.Scopes...)
	scopes = append(scopes, route.Scopes...)

	return scopes
}

func getRequestId(r *http.Request) string {
	requestID := uuid.NewString()

	if r == nil {
		return requestID
	}

	if xRequestIdHeader := r.Header.Get("X-RequestId"); len(xRequestIdHeader) > 0 {
		requestID = xRequestIdHeader
	} else if xApiGatekeeperRequestIdHeader := r.Header.Get("X-Api-Gatekeeper-RequestId"); len(xApiGatekeeperRequestIdHeader) > 0 {
		requestID = xApiGatekeeperRequestIdHeader
	}

	return requestID
}

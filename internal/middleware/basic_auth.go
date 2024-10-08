package middleware

import (
	"context"
	"net/http"

	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/model"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type BasicAuthService interface {
	AuthenticateToken(string) (model.User, error)

	Authorize(model.User, []string) error
}

type BackendRouteHandlerFunc = func(http.ResponseWriter, *http.Request, config.Backend, config.Route)

type BasicAuth struct {
	Service BasicAuthService
}

func (BasicAuth) GetAllScopes(backend config.Backend, route config.Route) []string {
	scopes := make([]string, len(backend.Scopes)+len(route.Scopes))
	scopes = append(scopes, backend.Scopes...)
	scopes = append(scopes, route.Scopes...)

	return scopes
}

func (a BasicAuth) GuardBackendRoute(w http.ResponseWriter, r *http.Request, backend config.Backend, route config.Route, next BackendRouteHandlerFunc) {
	if route.IsPublic {
		next(w, r, backend, route)
		return
	}

	user, err := a.Service.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
	}

	if err := a.Service.Authorize(user, a.GetAllScopes(backend, route)); err != nil {
		httputil.WriteForbidden(w)
	}

	ctx := context.WithValue(r.Context(), "userId", user.ID)

	next(w, r.WithContext(ctx), backend, route)
}

func (a BasicAuth) Guard(w http.ResponseWriter, r *http.Request, scopes []string, next http.HandlerFunc) {
	user, err := a.Service.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
	}

	if err := a.Service.Authorize(user, scopes); err != nil {
		httputil.WriteForbidden(w)
	}

	ctx := context.WithValue(r.Context(), "userId", user.ID)

	next(w, r.WithContext(ctx))
}

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

func (a BasicAuth) GuardBackendRoute(w http.ResponseWriter, r *http.Request, backend config.Backend, route config.Route, next BackendRouteHandlerFunc) {
	if route.IsPublic {
		next(w, r, backend, route)
		return
	}

	user, err := a.Service.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
	}

	if err := a.Service.Authorize(user, route.Scopes); err != nil {
		httputil.WriteForbidden(w)
	}

	ctx := context.WithValue(r.Context(), "userId", user.ID)

	next(w, r.WithContext(ctx), backend, route)
}

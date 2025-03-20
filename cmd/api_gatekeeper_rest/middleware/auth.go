package middleware

import (
	"net/http"

	"github.com/gustapinto/api-gatekeeper/internal/config"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type Auth struct {
	authService AuthService
}

func NewAuth(authService AuthService) Auth {
	return Auth{
		authService: authService,
	}
}

func (a Auth) GuardBackendRoute(
	w http.ResponseWriter,
	r *http.Request,
	backend config.Backend,
	route config.Route,
	next GuardBackendRouteNextFunc,
) {
	requestID := getRequestId(r)

	if route.IsPublic {
		next(w, r, backend, route)
		return
	}

	user, err := a.authService.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := a.authService.Authorize(user, mergeScopes(backend, route)); err != nil {
		httputil.WriteForbidden(w)
		return
	}

	ctx := withUserID(r.Context(), user.ID)
	ctx = withRequestID(ctx, requestID)

	next(w, r.WithContext(ctx), backend, route)
}

func (a Auth) GuardApplicationRoute(
	w http.ResponseWriter,
	r *http.Request,
	backend config.Backend,
	route config.Route,
	next GuardApplicationRouteNextFunc,
) {
	requestID := getRequestId(r)

	if route.IsPublic {
		next(w, r)
		return
	}

	user, err := a.authService.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := a.authService.Authorize(user, mergeScopes(backend, route)); err != nil {
		httputil.WriteForbidden(w)
		return
	}

	ctx := withUserID(r.Context(), user.ID)
	ctx = withRequestID(ctx, requestID)

	next(w, r.WithContext(ctx))
}

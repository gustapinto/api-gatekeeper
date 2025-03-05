package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/model"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type GuardBackendRouteNextFunc = func(http.ResponseWriter, *http.Request, config.Backend, config.Route)

type GuardApplicationRouteNextFunc = http.HandlerFunc

type BasicAuthService interface {
	AuthenticateToken(string) (model.User, error)

	Authorize(model.User, []string) error
}

type BasicAuth struct {
	basicAuthService BasicAuthService
}

func NewBasicAuth(basicAuthService BasicAuthService) BasicAuth {
	return BasicAuth{
		basicAuthService: basicAuthService,
	}
}

func (BasicAuth) getAllScopes(backend config.Backend, route config.Route) []string {
	scopes := make([]string, 0)
	scopes = append(scopes, backend.Scopes...)
	scopes = append(scopes, route.Scopes...)

	return scopes
}

func (BasicAuth) getRequestId(r *http.Request) string {
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

func (a BasicAuth) GuardBackendRoute(
	w http.ResponseWriter,
	r *http.Request,
	backend config.Backend,
	route config.Route,
	next GuardBackendRouteNextFunc,
) {
	requestID := a.getRequestId(r)

	if route.IsPublic {
		next(w, r, backend, route)
		return
	}

	user, err := a.basicAuthService.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := a.basicAuthService.Authorize(user, a.getAllScopes(backend, route)); err != nil {
		httputil.WriteForbidden(w)
		return
	}

	ctx := withUserID(r.Context(), user.ID)
	ctx = withRequestID(ctx, requestID)

	next(w, r.WithContext(ctx), backend, route)
}

func (a BasicAuth) GuardApplicationRoute(
	w http.ResponseWriter,
	r *http.Request,
	backend config.Backend,
	route config.Route,
	next GuardApplicationRouteNextFunc,
) {
	requestID := a.getRequestId(r)

	user, err := a.basicAuthService.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if err := a.basicAuthService.Authorize(user, a.getAllScopes(backend, route)); err != nil {
		httputil.WriteForbidden(w)
		return
	}

	ctx := withUserID(r.Context(), user.ID)
	ctx = withRequestID(ctx, requestID)

	next(w, r.WithContext(ctx))
}

package gateway

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gustapinto/api-gatekeeper/modules/user"
	"github.com/gustapinto/api-gatekeeper/util"
)

type AuthService interface {
	AuthenticateToken(string) (user.User, error)

	Authorize(user.User, []string) error
}

var alreadyRegisteredRoutes map[string]bool

func authorizeRoute(w http.ResponseWriter, r *http.Request, auth AuthService, route Route) (user.User, error) {
	if route.IsPublic {
		return user.User{}, nil
	}

	u, err := auth.AuthenticateToken(r.Header.Get("Authorization"))
	if err != nil {
		util.WriteUnauthorized(w)
		return user.User{}, err
	}

	if err := auth.Authorize(u, route.Scopes); err != nil {
		util.WriteForbidden(w)
		return user.User{}, err
	}

	return u, nil
}

func registerRoute(mux *http.ServeMux, logger *slog.Logger, auth AuthService, backend Backend, route Route) {
	routeLogger := logger.With("route", route.Name())
	routeIdentifier := route.QualifiedName(backend.Name)
	if _, exists := alreadyRegisteredRoutes[routeIdentifier]; exists {
		routeLogger.Warn("Route already registered, skipping")
		return
	}

	mux.HandleFunc(route.GatekeeperPath, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		user, err := authorizeRoute(w, r, auth, route)
		if err != nil {
			return
		}

		handleBackendRouteRequest(user.ID, backend, route, w, r)

		requestDuration := time.Since(start)

		routeLogger.Info("Request processed", "timeTaken", requestDuration)
	})

	if alreadyRegisteredRoutes == nil {
		alreadyRegisteredRoutes = make(map[string]bool)
	}

	alreadyRegisteredRoutes[routeIdentifier] = true
	routeLogger.Info("Route registered", "method", route.Method, "path", route.GatekeeperPath)
}

func RegisterModule(mux *http.ServeMux, logger *slog.Logger, auth AuthService, backends []Backend) {
	moduleLogger := logger.With("module", "gateway")

	for _, backend := range backends {
		backendLogger := moduleLogger.With("backend", backend.Name)

		for _, route := range backend.Routes {
			registerRoute(mux, backendLogger, auth, backend, route)
		}
	}
}

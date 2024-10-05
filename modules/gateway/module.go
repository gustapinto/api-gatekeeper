package gateway

import (
	"log/slog"
	"net/http"
	"time"
)

var alreadyRegisteredRoutes map[string]bool

func registerRoute(mux *http.ServeMux, logger *slog.Logger, backend Backend, route Route) {
	routeLogger := logger.With("route", route.Name())
	routeIdentifier := route.QualifiedName(backend.Name)
	if _, exists := alreadyRegisteredRoutes[routeIdentifier]; exists {
		routeLogger.Warn("Route already registered, skipping")
		return
	}

	mux.HandleFunc(route.GatekeeperPath, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userId := ""
		if !route.IsPublic {
			// TODO: Authenticate/Authorize user and update "userId"
		}

		handleBackendRouteRequest(userId, backend, route, w, r)

		requestDuration := time.Since(start)

		routeLogger.Info("Request processed", "timeTaken", requestDuration)
	})

	if alreadyRegisteredRoutes == nil {
		alreadyRegisteredRoutes = make(map[string]bool)
	}

	alreadyRegisteredRoutes[routeIdentifier] = true
	routeLogger.Info("Route registered", "method", route.Method, "path", route.GatekeeperPath)
}

func registerBackend(mux *http.ServeMux, logger *slog.Logger, backend Backend) {
	backendLogger := logger.With("backend", backend.Name)

	for _, route := range backend.Routes {
		registerRoute(mux, backendLogger, backend, route)
	}
}

func RegisterModule(mux *http.ServeMux, logger *slog.Logger, backends []Backend) {
	moduleLogger := logger.With("module", "gateway")

	for _, backend := range backends {
		registerBackend(mux, moduleLogger, backend)
	}
}

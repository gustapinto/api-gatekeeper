package main

import (
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/handler"
	"github.com/gustapinto/api-gatekeeper/internal/middleware"
	"github.com/gustapinto/api-gatekeeper/internal/repository/postgres"
	"github.com/gustapinto/api-gatekeeper/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	start := time.Now()

	configPath := flag.String("config", "", "The path to the config file")
	flag.Parse()

	config, err := config.LoadConfigFromYamlFile(configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if err := config.ValidateAndNormalize(); err != nil {
		logger.Error("Failed to validate config", "error", err)
		os.Exit(1)
	}

	logger.Info("Loaded application config from file", "configPath", *configPath)

	db, err := postgres.Conn{}.OpenDatabaseConnection(config.Database.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	userRepository := postgres.User{
		DB: db,
	}
	userService := service.User{
		Repository: userRepository,
	}
	basicAuth := middleware.BasicAuth{
		Service: userService,
	}
	backendHandler := handler.BackendHandler{
		Service: service.Backend{},
	}

	mux := http.NewServeMux()
	alreadyRegisteredRoutes := make(map[string]bool)
	for _, backend := range config.Backends {
		backendLogger := logger.With("backend", backend.Name)

		for _, route := range backend.Routes {
			routeLogger := backendLogger.With("route", route.Name())

			if _, exists := alreadyRegisteredRoutes[route.GatekeeperPath]; exists {
				routeLogger.Warn("Route already registered, skipping")
				return
			}

			mux.HandleFunc(route.GatekeeperPath, func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				basicAuth.GuardBackendRoute(w, r, backend, route, backendHandler.HandleBackendRouteRequest)

				requestDuration := time.Since(start)
				routeLogger.Info("Request processed", "timeTaken", requestDuration)
			})

			routeLogger.Info("Route registered", "method", route.Method, "path", route.GatekeeperPath)

			alreadyRegisteredRoutes[route.GatekeeperPath] = true
		}
	}

	address := config.API.Address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen", "address", address, "error", err.Error())
		os.Exit(1)
	}

	startupDuration := time.Since(start)

	logger.Info("Application started", "timeTaken", startupDuration, "address", address)

	err = http.Serve(listener, mux)
	if err != nil {
		logger.Error("Failed to server", "address", address, "error", err.Error())
		os.Exit(1)
	}
}

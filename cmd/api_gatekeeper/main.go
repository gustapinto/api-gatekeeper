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

	cfg, err := config.LoadConfigFromYamlFile(configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("Loaded application config from file", "configPath", *configPath)

	if err := cfg.ValidateAndNormalize(); err != nil {
		logger.Error("Failed to validate config", "error", err)
		os.Exit(1)
	}

	logger.Info("Validated application config")

	db, err := postgres.Conn{}.OpenDatabaseConnection(cfg.Database.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Connected to database")

	err = postgres.Conn{}.InitializeDatabase(db, cfg.API.User.Login, cfg.API.User.Password)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Initialized database data and application user")

	userRepository := postgres.User{DB: db}
	userService := service.User{Repository: userRepository}
	userHandler := handler.User{Service: userService}
	basicAuth := middleware.BasicAuth{Service: userService}
	backendHandler := handler.BackendHandler{Service: service.Backend{}}

	backends := cfg.Backends
	backends = append(backends, config.APIGatekeeperBackend(userHandler))

	logger.Info("Created dependencies")

	mux := http.NewServeMux()
	alreadyRegisteredRoutes := make(map[string]bool)
	for _, backend := range backends {
		backendLogger := logger.With("backend", backend.Name)

		for _, route := range backend.Routes {
			routeLogger := backendLogger.With("route", route.Name())

			if _, exists := alreadyRegisteredRoutes[route.Pattern()]; exists {
				routeLogger.Warn("Route already registered, skipping")
				return
			}

			mux.HandleFunc(route.Pattern(), func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				if route.HandlerFunc != nil {
					basicAuth.Guard(w, r, basicAuth.GetAllScopes(backend, route), route.HandlerFunc)
				} else {
					basicAuth.GuardBackendRoute(w, r, backend, route, backendHandler.HandleBackendRouteRequest)
				}

				requestDuration := time.Since(start)
				routeLogger.Info("Request processed", "timeTaken", requestDuration)
			})

			routeLogger.Info("Route registered", "method", route.Method, "path", route.GatekeeperPath)

			alreadyRegisteredRoutes[route.Pattern()] = true
		}
	}

	logger.Info("Registered all backends")

	address := cfg.API.Address
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

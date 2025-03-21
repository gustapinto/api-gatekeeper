package main

import (
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gustapinto/api-gatekeeper/cmd/api_gatekeeper_rest/handler"
	"github.com/gustapinto/api-gatekeeper/cmd/api_gatekeeper_rest/middleware"
	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/repository/gorm"
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

	db, err := gorm.OpenDatabaseConnection(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	logger.Info("Connected to database")

	userRepository := gorm.NewUser(db)
	basicAuthService := service.NewBasicAuth(userRepository)
	jwtService := service.NewJWT(userRepository, cfg.API.JwtSecret, cfg.API.TokenDuration())
	userService := service.NewUser(userRepository)
	userHandler := handler.NewUser(userService, jwtService)
	backendService := service.NewBackend()
	backendHandler := handler.NewBackend(backendService)

	backends := append(cfg.Backends, config.Backend{}.APIGatekeeperBackend(userHandler))

	logger.Info("Created dependencies")

	err = gorm.InitializeDatabase(db)
	if err != nil {
		logger.Error("Failed to initialize database schema", "error", err)
		os.Exit(1)
	}

	logger.Info("Initialized database schema")

	if err := userService.CreateApplicationUser(cfg.API.User); err != nil {
		logger.Error("Failed to initialize aplication user", "error", err)
		os.Exit(1)
	}

	logger.Info("Initialized application user")

	var authService middleware.AuthService
	switch cfg.API.AuthType {
	case config.AuthTypeBasic:
		authService = basicAuthService
	case config.AuthTypeJwt:
		authService = jwtService
	}

	auth := middleware.NewAuth(authService)

	mux := http.NewServeMux()
	alreadyRegisteredRoutes := make(map[string]bool)
	for _, backend := range backends {
		backendLogger := logger.With("backend", backend.Name)

		for _, route := range backend.Routes {
			routeLogger := backendLogger.With("route", route.Name())
			routePattern := route.Pattern()

			if _, exists := alreadyRegisteredRoutes[routePattern]; exists {
				routeLogger.Warn("Route already registered, skipping")
				continue
			}

			mux.HandleFunc(routePattern, func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				if route.IsApplicationRoute() {
					auth.GuardApplicationRoute(w, r, backend, route, route.HandlerFunc)
				} else {
					auth.GuardBackendRoute(w, r, backend, route, backendHandler.HandleBackendRouteRequest)
				}

				requestDuration := time.Since(start)
				routeLogger.Info("Request processed", "timeTaken", requestDuration)
			})

			routeLogger.Info("Route registered", "method", route.Method, "path", route.GatekeeperPath)

			alreadyRegisteredRoutes[routePattern] = true
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

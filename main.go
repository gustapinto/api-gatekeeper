package main

import (
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gustapinto/api-gatekeeper/modules/gateway"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	start := time.Now()

	configPath := flag.String("config", "", "The path to the config file")
	flag.Parse()

	config, err := LoadConfig(configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if err := config.ValidateAndNormalize(); err != nil {
		logger.Error("Failed to validate config", "error", err)
		os.Exit(1)
	}

	logger.Info("Loaded application config from file", "configPath", *configPath)

	mux := http.NewServeMux()

	gateway.RegisterModule(mux, logger, config.Backends)

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

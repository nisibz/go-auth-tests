package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nisibz/go-auth-tests/internal/adapter/config"
	"github.com/nisibz/go-auth-tests/internal/adapter/handler/http"
	"github.com/nisibz/go-auth-tests/internal/adapter/logger"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	logger.Set(config.App)

	router, err := http.NewRouter(
		config.HTTP,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}

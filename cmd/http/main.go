package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nisibz/go-auth-tests/internal/adapter/config"
	"github.com/nisibz/go-auth-tests/internal/adapter/handler/http"
	"github.com/nisibz/go-auth-tests/internal/adapter/logger"
	"github.com/nisibz/go-auth-tests/internal/adapter/storeages/mongodb"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	logger.Set(appConfig.App)

	slog.Info("Connecting to MongoDB...")
	mongoClient, err := mongodb.NewMongoClient(appConfig.Mongo)
	if err != nil {
		slog.Error("Error connecting to MongoDB", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to MongoDB")

	// Graceful shutdown for MongoDB client
	defer func() {
		slog.Info("Closing MongoDB connection...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			slog.Error("Error closing MongoDB connection", "error", err)
		}
		slog.Info("MongoDB connection closed.")
	}()

	router, err := http.NewRouter(
		appConfig.HTTP,
		// mongoClient, // You might pass the client to your router/handlers later
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting the application", "app", appConfig.App.Name, "env", appConfig.App.Env)

	listenAddr := fmt.Sprintf("%s:%s", appConfig.HTTP.URL, appConfig.HTTP.Port)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}

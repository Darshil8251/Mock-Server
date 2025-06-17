package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mock-server/internal/config"
	"mock-server/internal/server"
	"mock-server/pkg/logger"

	godotenv "github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Initialize logger
	Level := os.Getenv("LEVEL")
	mockLogger, err := logger.CreateNewLogger(Level)
	if err != nil {
		log.Fatal("Error initializing logger")
	}
	mockLogger.Info("logger successfully initialized")

	// Load config
	configFilePath := os.Getenv("CONFIG_PATH")
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		mockLogger.Error("Error loading config", err)
		return
	}
	mockLogger.InfoW("config successfully loaded", map[string]any{"config": config})

	// Create context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // Kubernetes/Docker stop
		syscall.SIGQUIT) // Graceful shutdown
	defer stop()

	PORT := os.Getenv("PORT")

	// Start the server
	if err := server.CreateServer(ctx, PORT, config); err != nil {
		mockLogger.Error("Server error", err)
		os.Exit(1)
	}

	// Block until context is cancelled (signal received)
	<-ctx.Done()
	mockLogger.Info("Received shutdown signal")
}
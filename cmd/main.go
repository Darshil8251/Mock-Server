package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"mock-server/internal/config"
	"mock-server/internal/server"
	"mock-server/pkg/logger"

	godotenv "github.com/joho/godotenv"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
		return
	}

	// Get environment variables
	configFilePath := os.Getenv("CONFIG_PATH")
	PORT := os.Getenv("PORT")
	Level := os.Getenv("LEVEL")



	// Initialize mockLogger
	mockLogger, err := logger.CreateNewLogger(Level)
	if err != nil {
		log.Fatal("Error initializing logger")
		return
	}

	mockLogger.Info("logger successfully initialized")

	// Load config
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		mockLogger.Error("Error loading config", err)
		return
	}

	mockLogger.InfoW("config successfully loaded", logger.Field("config", config))

	err = server.CreateServer(ctx, PORT, config)
	if err != nil {
		mockLogger.Error("Error starting server", err)
		return
	}

	select {
	case <-ctx.Done():
		mockLogger.Info("Server shutting down gracefully")
		return
	}
}

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

)

func main() {

	mockLogger, err := logger.CreateNewLogger()
	if err != nil {
		log.Fatal("Error initializing logger")
	}
	mockLogger.Info("logger successfully initialized")

	config, err := config.LoadConfig()
	if err != nil {
		mockLogger.Error("Error loading config", err)
		return
	}
	mockLogger.InfoW("config successfully loaded", map[string]any{"config": config})

	
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // Kubernetes/Docker stop
		syscall.SIGQUIT) // Graceful shutdown
	defer stop()

	if err := server.CreateServer(ctx, config); err != nil {
		mockLogger.Error("Server error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	mockLogger.Info("Received shutdown signal")
}

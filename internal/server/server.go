package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"mock-server/internal/config"
	"mock-server/internal/middleware"
	"mock-server/internal/router"
	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CreateServer(ctx context.Context, port string, cfg *config.APIConfig) error {
	tmpLogger := logger.GetLogger()

	// Create a new Gin engine
	engine := gin.New()

	// Recovery should come first
	engine.Use(gin.Recovery())

	// Then your logging middleware
	engine.Use(middleware.LoggerMiddleware())

	// Initialize the router AFTER middleware
	err := router.SetupRoutes(engine, cfg)
	if err != nil {
		return err
	}

	// Create HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Channel to listen for server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		tmpLogger.InfoW("Server starting", map[string]any{"port": port})
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server failed: %w", err)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		tmpLogger.Info("Server shutting down gracefully...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		tmpLogger.Info("Server stopped")
		return nil
	}
}

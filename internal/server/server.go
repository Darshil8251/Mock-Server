package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"mock-server/internal/config"
	"mock-server/internal/middleware"
	"mock-server/internal/router"
	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CreateServer(ctx context.Context, cfg *config.APIConfig) error {
	var (
		mockLogger = logger.GetLogger()
		port      = os.Getenv("PORT")
	)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.LoggerMiddleware())

	// Setup routers
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

	serverErr := make(chan error, 1)

	go func() {
		mockLogger.InfoW("Server starting", map[string]any{"port": port})
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server failed: %w", err)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		mockLogger.Info("Server shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		mockLogger.Info("Server stopped")
		return nil
	}
}

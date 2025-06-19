package server

import (
	"context"
	"fmt"

	"mock-server/internal/config"
	"mock-server/internal/router"
	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CreateServer(ctx context.Context, port string, cfg *config.APIConfig) (err error) {

	tmpLogger := logger.GetLogger()

	// Create a new Gin engine
	engine := gin.Default()

	// Initialize the router with the provided configuration
	router.SetupRoutes(engine, cfg)

	// Start the server with the specified configuration
	if err := engine.Run(":" + port); err != nil {
		serverErr := fmt.Errorf("failed to start server on port %s: %w", port, err)
		tmpLogger.Error("Failed to start server", serverErr)
		return serverErr
	}

	return nil
}

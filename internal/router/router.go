package router

import (
	"mock-server/internal/config"
	"mock-server/internal/pagination"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes based on the API config
func SetupRoutes(engine *gin.Engine, cfg *config.APIConfig) {

	// Register all endpoints from config
	for _, endpoint := range cfg.Endpoints {

		// Register the handler with the appropriate HTTP method
		switch endpoint.Method {
		case "GET":
			engine.GET(endpoint.Path, pagination.CreatePaginator(endpoint))
		case "POST":
			engine.POST(endpoint.Path, pagination.CreatePaginator(endpoint))
		case "PUT":
			engine.PUT(endpoint.Path, pagination.CreatePaginator(endpoint))
		case "DELETE":
			engine.DELETE(endpoint.Path, pagination.CreatePaginator(endpoint))
		case "PATCH":
			engine.PATCH(endpoint.Path, pagination.CreatePaginator(endpoint))
		default:
			engine.Any(endpoint.Path, pagination.CreatePaginator(endpoint))
		}
	}
}

package router

import (
	"errors"
	
	"mock-server/internal/config"
	"mock-server/internal/pagination"

	"github.com/gin-gonic/gin"
)

var (
	errSetupRoutes = errors.New("failed to setup routes")
)

// SetupRoutes configures all routes based on the API config
func SetupRoutes(engine *gin.Engine, cfg *config.APIConfig) error {

	// Register all endpoints from config
	for _, endpoint := range cfg.Endpoints {

		// Create the paginator for the endpoint
		paginator, err := pagination.CreatePaginator(endpoint)
		if err != nil {
			return errors.Join(errSetupRoutes, err)
		}

		// Register the handler with the appropriate HTTP method
		switch endpoint.Method {
		case "GET":
			engine.GET(endpoint.Path, paginator.Paginate)
		case "POST":
			engine.POST(endpoint.Path, paginator.Paginate)
		case "PUT":
			engine.PUT(endpoint.Path, paginator.Paginate)
		case "DELETE":
			engine.DELETE(endpoint.Path, paginator.Paginate)
		case "PATCH":
			engine.PATCH(endpoint.Path, paginator.Paginate)
		default:
			engine.Any(endpoint.Path, paginator.Paginate)
		}
	}
	return nil
}

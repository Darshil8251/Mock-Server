package server

import (
	"strconv"

	"mock-server/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *models.Config
	router *gin.Engine
}

func NewServer(config *models.Config) *Server {
	server := &Server{
		config: config,
		router: gin.Default(),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	for _, endpoint := range s.config.Endpoints {
		s.router.Handle(endpoint.Method, endpoint.Path, s.handleEndpoint(endpoint))
	}
}

func (s *Server) handleEndpoint(endpoint models.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set response headers
		for key, value := range endpoint.Response.Headers {
			c.Header(key, value)
		}

		// Handle pagination
		if endpoint.Pagination.Enabled {
			page, _ := strconv.Atoi(c.DefaultQuery(endpoint.Pagination.PageParam, "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery(endpoint.Pagination.SizeParam, strconv.Itoa(endpoint.Pagination.PageSize)))

			// Calculate pagination metadata
			totalPages := (endpoint.Pagination.TotalItems + pageSize - 1) / pageSize
			if page > totalPages {
				page = totalPages
			}

			// Modify response based on pagination
			response := modifyResponseForPagination(endpoint.Response.Body, page, pageSize, endpoint.ResponseParams)

			c.JSON(endpoint.Response.Status, gin.H{
				"data": response,
				"pagination": gin.H{
					"currentPage": page,
					"pageSize":    pageSize,
					"totalPages":  totalPages,
					"totalItems":  endpoint.Pagination.TotalItems,
				},
			})
			return
		}

		// Return regular response without pagination
		c.JSON(endpoint.Response.Status, endpoint.Response.Body)
	}
}

func modifyResponseForPagination(response interface{}, page, pageSize int, params models.ResponseParams) interface{} {
	// This is a simplified version. In a real implementation, you would need to
	// handle different types of responses and modify them according to the pagination
	// and dynamic fields configuration.
	return response
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

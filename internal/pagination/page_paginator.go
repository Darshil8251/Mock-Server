package pagination

import (
	"mock-server/internal/config"

	"github.com/gin-gonic/gin"
)

func pagePaginator(endpoint config.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example handler logic
		c.JSON(200, gin.H{
			"message": "This is a page paginator",
		})
	}
}

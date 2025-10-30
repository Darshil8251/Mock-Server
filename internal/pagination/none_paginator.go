package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"mock-server/internal/config"
	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// pagePaginator responsible for the page based pagination
type defaultPagination struct {
	responseObj map[string]any
}

// createPagePaginator creates a new page paginator for the given endpoint
func createDefaultPaginator(endpoint config.Endpoint) (*defaultPagination, error) {
	var mockLogger = logger.GetLogger().With(map[string]any{"endpoint": endpoint.Path})

	mockLogger.Info("creating  default pagination")

	p := &defaultPagination{}

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		mockLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	return p, nil
}

// Paginate is the handler function for the page paginator
func (p *defaultPagination) Paginate(c *gin.Context) {

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

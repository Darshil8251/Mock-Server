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
	responseObj       map[string]any
	pageRecordCount   int
	totalRecordCount  int
	responseFieldName string
	sendRecordCount   int
}

// createPagePaginator creates a new page paginator for the given endpoint
func createDefaultPaginator(endpoint config.Endpoint) (*defaultPagination, error) {
	var mockLogger = logger.GetLogger()

	mockLogger.InfoW("creating the default pagination", map[string]any{"endpoint": endpoint.Path})

	p := &defaultPagination{}

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		mockLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	options := endpoint.Pagination.Options

	p.pageRecordCount = 100
	p.totalRecordCount = 1000
	p.sendRecordCount = 0

	if count, ok := options["pageRecordCount"].(int); ok {
		p.pageRecordCount = count

	}

	if count, ok := options["totalRecordCount"].(int); ok {
		p.totalRecordCount = count

	}
	p.responseFieldName, err = findResponseFieldName(endpoint.Response.FieldName, p.responseObj)
	if err != nil {
		return nil, errors.Join(errors.New("error to find response field"), err)
	}

	return p, nil
}

// Paginate is the handler function for the page paginator
func (p *defaultPagination) Paginate(c *gin.Context) {

	if p.sendRecordCount >= p.totalRecordCount {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
	}

	arr, ok := p.responseObj[p.responseFieldName].([]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, p.pageRecordCount)

	for len(APIResponseObject) < p.pageRecordCount {
		APIResponseObject = append(APIResponseObject, object)
	}

	p.sendRecordCount += len(APIResponseObject)

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

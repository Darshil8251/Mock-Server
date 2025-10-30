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
type tokenPaginator struct {
	responseObj          map[string]interface{}
	tokenFieldName       string
	responseFieldName    string
	paginationParameters paginationParameters
}

// createPagePaginator creates a new page paginator for the given endpoint
func createTokenPaginator(endpoint config.Endpoint) (*tokenPaginator, error) {
	var mockLogger = logger.GetLogger().With(map[string]any{"endpoint": endpoint.Path})

	mockLogger.Info("creating token paginator")

	t := &tokenPaginator{}
	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		mockLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	t.responseObj = responseObj

	t.paginationParameters = loadPaginationParameters(endpoint)

	tokenField, ok := endpoint.Pagination.Options["tokenFieldName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token field name for endpoint: %s", endpoint.Path)
	}
	t.tokenFieldName = tokenField
	t.responseFieldName = endpoint.Response.FieldName

	return t, nil
}

// Paginate is the handler function for the page paginator
func (t *tokenPaginator) Paginate(c *gin.Context) {

	if t.paginationParameters.sendRecordsCount >= t.paginationParameters.totalRecord {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No record found"})
		return
	}

	// Get single object from response field
	arr, ok := t.responseObj[t.responseFieldName].([]any)
	if !ok || len(arr) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}
	object := arr[0]

	numItems := t.paginationParameters.pageSize

	// last page
	if t.paginationParameters.sendPageCount+1 == t.paginationParameters.totalPageCount {
		numItems = t.paginationParameters.totalRecord - t.paginationParameters.sendRecordsCount
	}

	// Build response array
	APIResponseObject := make([]any, 0, numItems)
	for i := 0; i < numItems; i++ {
		APIResponseObject = append(APIResponseObject, object)
	}
	t.responseObj[t.responseFieldName] = APIResponseObject
	t.paginationParameters.sendRecordsCount += numItems
	t.paginationParameters.sendPageCount++

	// Set next offset value
	t.responseObj[t.tokenFieldName] = "abcdxyzed"

	if t.paginationParameters.sendRecordsCount >= t.paginationParameters.totalRecord {
		t.responseObj[t.tokenFieldName] = nil
	}

	jsonResponse, err := json.Marshal(t.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

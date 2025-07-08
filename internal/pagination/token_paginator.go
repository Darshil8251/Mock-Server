package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"mock-server/internal/config"
	"mock-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

// pagePaginator responsible for the page based pagination
type tokenPaginator struct {
	responseObj           map[string]interface{}
	tokenLocation         pageParameterLocation
	responseField         string
	paginationParameters  paginationParameters
}

// createPagePaginator creates a new page paginator for the given endpoint
func createTokenPaginator(endpoint config.Endpoint) (*tokenPaginator, error) {
	var mockLogger = logger.GetLogger()

	mockLogger.InfoW("creating page paginator", map[string]any{"endpoint": endpoint.Path})

	t := tokenPaginator{
		tokenLocation: pageParameterLocation(endpoint.Pagination.Location),
	}
	responseObj, err := loadResponseObj(endpoint.ResponseObjFilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		mockLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	t.responseObj = responseObj

	t.paginationParameters = loadPaginationParameters(endpoint)

	// Validate the response field
	if endpoint.ResponseField != "" {
		_, ok := responseObj[endpoint.ResponseField].([]any)
		if !ok {
			errInvalidResponseField := fmt.Errorf("invalid response field for endpoint: %v", endpoint.Path)
			mockLogger.Warn(errInvalidResponseField.Error(), err)
			return nil, errors.Join(errInvalidResponseField, err)
		}
		t.responseField = endpoint.ResponseField
		return &t, nil
	}

	if endpoint.ResponseField == "" {
		for k, v := range responseObj {
			if _, ok := v.([]interface{}); ok {
				t.responseField = k
				return &t, nil
			}
		}
	}

	errInvalidResponseField := fmt.Errorf("response field not present in response object for endpoint: %v", endpoint.Path)
	mockLogger.Warn(errInvalidResponseField.Error(), err)
	return nil, errors.Join(errInvalidResponseField, err)
}

// Paginate is the handler function for the page paginator
func (t *tokenPaginator) Paginate(c *gin.Context) {

	var pageSize = defaultPageSize

	if t.paginationParameters.pageSentCount >= t.paginationParameters.totalPageCount {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
	}

	// Extract pagination params from the respective location
	switch t.tokenLocation {
	case body:
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
			return
		}
		value, found := requestBody[t.paginationParameters.pageSizeKey]
		if found {
			switch v := value.(type) {
			case float64:
				pageSize = int(v)
			case int:
				// Already an integer
				pageSize = v
			case int32, int64:
				// Handle other integer types
				pageSize = int(reflect.ValueOf(v).Int())
			case uint, uint32, uint64:
				// Handle unsigned integers
				pageSize = int(reflect.ValueOf(v).Uint())
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "size must be a number"})
				return
			}
		}

	case header:
		if v := c.GetHeader(t.paginationParameters.pageSizeKey); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				pageSize = p
			}
		}
	case query:
		size, err := strconv.Atoi(c.DefaultQuery(t.paginationParameters.pageSizeKey, strconv.Itoa(defaultPageSize)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sizeValue"})
			return
		}
		pageSize = size
	}

	// 3. Find the response object
	arr, ok := t.responseObj[t.responseField].([]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	if t.paginationParameters.sentRecordsCount+pageSize > t.paginationParameters.totalRecordCount {
		pageSize = t.paginationParameters.totalRecordCount - t.paginationParameters.sentRecordsCount
	}

	for len(APIResponseObject) < pageSize {
		APIResponseObject = append(APIResponseObject, object)
	}

	t.paginationParameters.pageSentCount++
	t.paginationParameters.sentRecordsCount += pageSize

	t.responseObj[t.responseField] = APIResponseObject

	jsonResponse, err := json.Marshal(t.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

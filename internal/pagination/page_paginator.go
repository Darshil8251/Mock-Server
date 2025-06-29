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
type pagePaginator struct {
	responseObj          map[string]interface{}
	pageParamsLocation   pageParameterLocation
	responseField        string
	paginationParameters paginationParameters
}

// createPagePaginator creates a new page paginator for the given endpoint
func createPagePaginator(endpoint config.Endpoint) (*pagePaginator, error) {
	var mockLogger = logger.GetLogger()

	mockLogger.InfoW("creating page paginator", map[string]any{"endpoint": endpoint.Path})

	p := pagePaginator{
		pageParamsLocation: pageParameterLocation(endpoint.Pagination.Location),
	}

	responseObj, err := loadResponseObj(endpoint.ResponseObjFilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		mockLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	p.paginationParameters = loadPaginationParameters(endpoint)

	// Validate the response field
	if endpoint.ResponseField != "" {
		_, ok := responseObj[endpoint.ResponseField].([]any)
		if !ok {
			errInvalidResponseField := fmt.Errorf("invalid response field for endpoint: %v", endpoint.Path)
			mockLogger.Warn(errInvalidResponseField.Error(), err)
			return nil, errors.Join(errInvalidResponseField, err)
		}
		p.responseField = endpoint.ResponseField
		return &p, nil
	}

	if endpoint.ResponseField == "" {
		for k, v := range responseObj {
			if _, ok := v.([]interface{}); ok {
				p.responseField = k
				return &p, nil
			}
		}
	}

	errInvalidResponseField := fmt.Errorf("response field not present in response object for endpoint: %v", endpoint.Path)
	mockLogger.Warn(errInvalidResponseField.Error(), err)
	return nil, errors.Join(errInvalidResponseField, err)
}

// Paginate is the handler function for the page paginator
func (p *pagePaginator) Paginate(c *gin.Context) {

	var pageSize = defaultPageSize

	if p.paginationParameters.pageSentCount >= p.paginationParameters.totalPageCount {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
	}

	// Extract pagination params from the respective location
	switch p.pageParamsLocation {
	case body:
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
			return
		}
		value, found := requestBody[p.paginationParameters.pageSizeKey]
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
		if v := c.GetHeader(p.paginationParameters.pageSizeKey); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				pageSize = p
			}
		}
	case query:
		size, err := strconv.Atoi(c.DefaultQuery(p.paginationParameters.pageSizeKey, strconv.Itoa(defaultPageSize)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sizeValue"})
			return
		}
		pageSize = size
	}

	// 3. Find the response object
	arr, ok := p.responseObj[p.responseField].([]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	if p.paginationParameters.sentRecordsCount+pageSize > p.paginationParameters.totalRecordCount {
		pageSize = p.paginationParameters.totalRecordCount - p.paginationParameters.sentRecordsCount
	}

	for len(APIResponseObject) < pageSize {
		APIResponseObject = append(APIResponseObject, object)
	}

	p.paginationParameters.pageSentCount++
	p.paginationParameters.sentRecordsCount += pageSize

	p.responseObj[p.responseField] = APIResponseObject

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

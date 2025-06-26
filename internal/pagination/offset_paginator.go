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

type offsetPaginator struct {
	responseObj          map[string]interface{}
	offsetKey            string
	limitKey             string
	offsetLocation       pageParameterLocation
	responseField        string
	paginationParameters paginationParameters
}

var _ Paginator = (*offsetPaginator)(nil)

func createOffsetPaginator(endpoint config.Endpoint) (Paginator, error) {
	var tmpLogger = logger.GetLogger()

	tmpLogger.InfoW("creating offset paginator", map[string]any{"endpoint": endpoint.Path})

	o := offsetPaginator{
		offsetLocation: pageParameterLocation(endpoint.Pagination.Location),
	}

	responseObj, err := loadResponseObj(endpoint.ResponseObjFilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		tmpLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	o.responseObj = responseObj

	o.paginationParameters = loadPaginationParameters(endpoint)

	// Validate the response field
	if endpoint.ResponseField != "" {
		_, ok := responseObj[endpoint.ResponseField].([]any)
		if !ok {
			errInvalidResponseField := fmt.Errorf("invalid response field for endpoint: %v", endpoint.Path)
			tmpLogger.Warn(errInvalidResponseField.Error(), err)
			return nil, errors.Join(errInvalidResponseField, err)
		}
		o.responseField = endpoint.ResponseField
		return &o, nil
	}

	// If user not specified the response field, then find array field from the response object
	if endpoint.ResponseField == "" {
		for k, v := range responseObj {
			if _, ok := v.([]interface{}); ok {
				o.responseField = k
				return &o, nil
			}
		}
	}

	errInvalidResponseField := fmt.Errorf("response field not present in response object for endpoint: %v", endpoint.Path)
	tmpLogger.Warn(errInvalidResponseField.Error(), err)
	return nil, errors.Join(errInvalidResponseField, err)
}

func (o *offsetPaginator) Paginate(c *gin.Context) {
	var (
		pageSize  = defaultPageSize
		tmpLogger = logger.GetLogger()
	)

	tmpLogger.InfoW("Paginate", map[string]any{
		"offsetKey": o.offsetKey,
		"limitKey":  o.limitKey,
		"location":  o.offsetLocation,
	})

	// Extract pagination params from the respective location
	switch o.offsetLocation {
	case body:
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
			return
		}
		value, found := requestBody[o.limitKey]
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
		if v := c.GetHeader(o.limitKey); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				pageSize = p
			}
		}
	case query:
		size, err := strconv.Atoi(c.DefaultQuery(o.limitKey, strconv.Itoa(defaultPageSize)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sizeValue"})
			return
		}
		pageSize = size
	}

	tmpLogger.InfoW("page value size", map[string]any{"size": pageSize})

	// 3. Find the response object
	arr, ok := o.responseObj[o.responseField].([]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	if o.paginationParameters.sentRecordsCount+pageSize > o.paginationParameters.totalRecordCount {
		pageSize = o.paginationParameters.totalRecordCount - o.paginationParameters.sentRecordsCount
	}

	for len(APIResponseObject) < int(pageSize) {
		APIResponseObject = append(APIResponseObject, object)
	}

	o.responseObj[o.responseField] = APIResponseObject

	jsonResponse, err := json.Marshal(o.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)
}

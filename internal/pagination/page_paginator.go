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
	"mock-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageKey  = "page"
	defaultPageSize = 100
)

// pagePaginator responsible for the page based pagination
type pagePaginator struct {
	endpoint           config.Endpoint
	responseObj        map[string]interface{}
	pageSizeKey        string
	pageKey            string
	pageParamsLocation pageParamsLocation
	responseField      string
}

// createPagePaginator creates a new page paginator for the given endpoint
func createPagePaginator(endpoint config.Endpoint) (*pagePaginator, error) {
	var tmpLogger = logger.GetLogger()

	p := &pagePaginator{
		endpoint: endpoint,
	}

	filePath, err := utils.ValidateAndResolvePath(endpoint.ResponseObjFilePath)
	if err != nil {
		errInvalidPath := fmt.Errorf("invalid file path for endpoint: %s", endpoint.Path)
		tmpLogger.Warn(errInvalidPath.Error(), err)
		return nil, errors.Join(errInvalidPath, err)
	}

	responseObj, err := utils.ReadAndParseJSONFile(filePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response object for endpoint: %s", endpoint.Path)
		tmpLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	// Validate the pagination parameters and set the values
	pageSizeKey, ok := endpoint.Pagination.Options["pageSizeKey"].(string)
	if !ok {
		errInvalidPageSizeKey := fmt.Errorf("invalid pagesize key for endpoint: %v", endpoint.Path)
		tmpLogger.Warn(errInvalidPageSizeKey.Error(), err)
		return nil, errors.Join(errInvalidPageSizeKey, err)
	}
	p.pageSizeKey = pageSizeKey

	pageKey, ok := endpoint.Pagination.Options["pageKey"].(string)
	if !ok {
		errInvalidPageKey := fmt.Errorf("invalid page key for endpoint: %v", endpoint.Path)
		tmpLogger.Warn(errInvalidPageKey.Error(), err)
		return nil, errors.Join(errInvalidPageKey, err)
	}
	p.pageKey = pageKey

	p.pageParamsLocation = pageParamsLocation(endpoint.Pagination.Location)

	// Set the response field name
	var arrayField = endpoint.ResponseField

	// If user not specified the response field, then find array field from the response object
	if endpoint.ResponseField == "" {
		for k, v := range responseObj {
			if _, ok := v.([]interface{}); ok {
				arrayField = k
				break
			}
		}
	}

	if arrayField == "" {
		errInvalidResponseField := fmt.Errorf("invalid response field for endpoint: %v", endpoint.Path)
		tmpLogger.Warn(errInvalidResponseField.Error(), err)
		return nil, errors.Join(errInvalidResponseField, err)
	}

	p.responseField = arrayField

	return p, nil
}

// Paginate is the handler function for the page paginator
func (p *pagePaginator) Paginate(c *gin.Context) {

	var pageSize = defaultPageSize

	// Extract pagination params from the respective location
	switch p.pageParamsLocation {
	case body:
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
			return
		}
		value, found := requestBody[p.pageSizeKey]
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
		if v := c.GetHeader(p.pageSizeKey); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				pageSize = p
			}
		}
	case query:
		size, err := strconv.Atoi(c.DefaultQuery(p.pageSizeKey, strconv.Itoa(defaultPageSize)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sizeValue"})
			return
		}
		pageSize = size
	}

	// 3. Find the response object
	arr, ok := p.responseObj[p.responseField].([]any)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find response object"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	for len(APIResponseObject) < int(pageSize) {
		APIResponseObject = append(APIResponseObject, object)
	}

	p.responseObj[p.responseField] = APIResponseObject

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"mock-server/internal/config"

	"github.com/gin-gonic/gin"
)

// pagePaginator responsible for the page based pagination
type pagePaginator struct {
	responseObj        map[string]any
	pageParamsLocation pageParameterLocation
	responseField      string
	pageSizeKey        string
	totalRecord        int
	sendRecordCount    int
}

// createPagePaginator creates a new page paginator for the given endpoint
func createPagePaginator(endpoint config.Endpoint) (*pagePaginator, error) {

	p := &pagePaginator{
		pageParamsLocation: pageParameterLocation(endpoint.Pagination.Location),
	}

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	totalRecord := defaultTotalRecordCount

	switch v := endpoint.Pagination.Options["totalRecord"].(type) {
	case int:
		totalRecord = v
	case float64:
		totalRecord = int(v)
	case int64:
		totalRecord = int(v)
	case int32:
		totalRecord = int(v)
	}

	p.totalRecord = totalRecord

	fmt.Println("total record count", p.totalRecord)
	p.responseField, err = findResponseFieldName(endpoint.Response.FieldName, p.responseObj)
	if err != nil {
		return nil, errors.Join(errors.New("error to find response field"), err)
	}

	return p, nil

}

// Paginate is the handler function for the page paginator
func (p *pagePaginator) Paginate(c *gin.Context) {
	if p.sendRecordCount >= p.totalRecord {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var pageSize = defaultPageSize
	// Extract pagination params from the respective location
	switch p.pageParamsLocation {
	case body:
		var requestBody map[string]any
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		value, found := requestBody[p.pageSizeKey]
		if found {
			switch v := value.(type) {
			case float64:
				pageSize = int(v)
			case int:
				pageSize = v
			case int32, int64:
				pageSize = int(reflect.ValueOf(v).Int())
			case uint, uint32, uint64:
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	if p.totalRecord-p.sendRecordCount < pageSize {

		pageSize = p.totalRecord - p.sendRecordCount
	}

	for len(APIResponseObject) < pageSize {
		APIResponseObject = append(APIResponseObject, object)
	}

	p.sendRecordCount += len(APIResponseObject)

	p.responseObj[p.responseField] = APIResponseObject

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

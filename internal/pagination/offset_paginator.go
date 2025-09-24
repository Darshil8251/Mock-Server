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
	offsetLocation       pageParameterLocation
	responseField        string
	paginationParameters paginationParameters
}

var _ Paginator = (*offsetPaginator)(nil)

func createOffsetPaginator(endpoint config.Endpoint) (Paginator, error) {
	var tmpLogger = logger.GetLogger()

	tmpLogger.InfoW("creating offset paginator", map[string]any{"endpoint": endpoint.Path})

	o := &offsetPaginator{
		offsetLocation: pageParameterLocation(endpoint.Pagination.Location),
	}

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		tmpLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	o.responseObj = responseObj

	o.paginationParameters = loadPaginationParameters(endpoint)

	o.responseField, err = findResponseFieldName(endpoint.Response.FieldName, o.responseObj)
	if err != nil {
		return nil, errors.Join(errors.New("error to find response field"), err)
	}

	return o, nil

}

func (o *offsetPaginator) Paginate(c *gin.Context) {
	var (
		pageSize  = defaultPageSize
		tmpLogger = logger.GetLogger()
	)

	// Extract pagination params from the respective location
	switch o.offsetLocation {
	case body:
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
			return
		}
		value, found := requestBody[o.paginationParameters.pageSizeKey]
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
		if v := c.GetHeader(o.paginationParameters.pageSizeKey); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				pageSize = p
			}
		}
	case query:
		size, err := strconv.Atoi(c.DefaultQuery(o.paginationParameters.pageSizeKey, strconv.Itoa(defaultPageSize)))
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

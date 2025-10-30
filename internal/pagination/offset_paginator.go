package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"mock-server/internal/config"
)

type offsetPaginator struct {
	responseObj          map[string]interface{}
	offsetKey            string
	responseField        string
	paginationParameters paginationParameters
}

var _ Paginator = (*offsetPaginator)(nil)

func createOffsetPaginator(endpoint config.Endpoint) (Paginator, error) {

	o := &offsetPaginator{}

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		return nil, errors.Join(errInvalidResponse, err)
	}

	o.responseObj = responseObj

	o.paginationParameters = loadPaginationParameters(endpoint)

	o.responseField = endpoint.Response.FieldName

	offsetField, ok := endpoint.Pagination.Options["offsetKey"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid offset field for the endpoint: %v", endpoint.Path)
	}
	o.offsetKey = offsetField

	return o, nil

}

func (o *offsetPaginator) Paginate(c *gin.Context) {
	if o.paginationParameters.sendRecordsCount >= o.paginationParameters.totalRecord {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No record found"})
		return
	}

	// Get single object from response field
	arr, ok := o.responseObj[o.responseField].([]any)
	if !ok || len(arr) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}
	object := arr[0]

	numItems := o.paginationParameters.pageSize

	// last page
	if o.paginationParameters.sendPageCount+1 == o.paginationParameters.totalPageCount {
		numItems = o.paginationParameters.totalRecord - o.paginationParameters.sendRecordsCount
	}

	// Build response array
	APIResponseObject := make([]any, 0, numItems)
	for i := 0; i < numItems; i++ {
		APIResponseObject = append(APIResponseObject, object)
	}
	o.responseObj[o.responseField] = APIResponseObject
	o.paginationParameters.sendRecordsCount += numItems
	o.paginationParameters.sendPageCount++

	// Set next offset value
	nextOffset := o.paginationParameters.sendRecordsCount
	o.responseObj[o.offsetKey] = nextOffset + 1

	if o.paginationParameters.sendRecordsCount >= o.paginationParameters.totalRecord {
		o.responseObj[o.offsetKey] = nil
	}

	jsonResponse, err := json.Marshal(o.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)
}

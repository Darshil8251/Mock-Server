package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"mock-server/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// pagePaginator responsible for the page based pagination
type pagePaginator struct {
	responseObj      map[string]any
	responseField    string
	paginationParams paginationParameters
}

// createPagePaginator creates a new page paginator for the given endpoint
func createPagePaginator(endpoint config.Endpoint) (*pagePaginator, error) {

	p := &pagePaginator{}

	p.paginationParams = loadPaginationParameters(endpoint)

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		return nil, errors.Join(errInvalidResponse, err)
	}

	p.responseObj = responseObj

	p.responseField = endpoint.Response.FieldName

	return p, nil

}

// Paginate is the handler function for the page paginator
func (p *pagePaginator) Paginate(c *gin.Context) {
	if p.paginationParams.sendRecordsCount >= p.paginationParams.totalRecord {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No record found"})
		return
	}

	// Get single object from response field
	arr, ok := p.responseObj[p.responseField].([]any)
	if !ok || len(arr) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}
	object := arr[0]

	numItems := p.paginationParams.pageSize

	// last page
	if p.paginationParams.sendPageCount+1 == p.paginationParams.totalPageCount {
		numItems = p.paginationParams.totalRecord - p.paginationParams.sendRecordsCount
	}

	// Build response array
	APIResponseObject := make([]any, 0, numItems)
	for i := 0; i < numItems; i++ {
		APIResponseObject = append(APIResponseObject, object)
	}
	p.responseObj[p.responseField] = APIResponseObject
	p.paginationParams.sendRecordsCount += numItems
	p.paginationParams.sendPageCount++

	jsonResponse, err := json.Marshal(p.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)
}

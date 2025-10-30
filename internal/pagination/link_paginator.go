package pagination

import (
	"encoding/json"
	"fmt"
	"mock-server/internal/config"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type linkPaginator struct {
	responseObj          map[string]interface{}
	linkKey              string
	responseField        string
	paginationParameters paginationParameters
}

var _ Paginator = (*linkPaginator)(nil)

func createLinkPaginator(endpoint config.Endpoint) (Paginator, error) {

	l := &linkPaginator{}

	l.paginationParameters = loadPaginationParameters(endpoint)
	l.linkKey = "nextPage"
	linkKeyName, ok := endpoint.Pagination.Options["linkKey"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid link key for the endpoint: %v", endpoint.Path)
	}
	l.linkKey = linkKeyName

	responseObj, err := loadResponseObj(endpoint.Response.FilePath)
	if err != nil {
		return nil, fmt.Errorf("invalid response file path for endpoint: %s, %w", endpoint.Path, err)
	}

	l.responseObj = responseObj

	l.responseField = endpoint.Response.FieldName

	return l, nil
}

func (l *linkPaginator) Paginate(c *gin.Context) {
	// Get single object from response field
	arr, ok := l.responseObj[l.responseField].([]any)
	if !ok || len(arr) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response field"})
		return
	}
	object := arr[0]

	// Calculate number of items for this page
	if l.paginationParameters.sendRecordsCount >= l.paginationParameters.totalRecord {
		c.JSON(http.StatusNotFound, gin.H{"msg": "No record found"})
		return
	}

	numItems := l.paginationParameters.pageSize

	// last page
	if l.paginationParameters.sendPageCount+1 == l.paginationParameters.totalPageCount {
		numItems = l.paginationParameters.totalRecord - l.paginationParameters.sendRecordsCount
	}

	// Build response array
	APIResponseObject := make([]any, 0, numItems)
	for i := 0; i < numItems; i++ {
		APIResponseObject = append(APIResponseObject, object)
	}
	l.responseObj[l.responseField] = APIResponseObject
	l.paginationParameters.sendRecordsCount += numItems
	l.paginationParameters.sendPageCount++

	// Generate next page link or set to null
	if l.paginationParameters.sendPageCount < l.paginationParameters.totalPageCount {
		l.responseObj[l.linkKey] = generatePageLink(c, l.paginationParameters.sendPageCount+1, l.paginationParameters.pageSize)
	} else {
		l.responseObj[l.linkKey] = nil
	}

	jsonResponse, err := json.Marshal(l.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonResponse)
}

func generatePageLink(c *gin.Context, nextPage int, pageSize int) string {
	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	u := &url.URL{
		Scheme: scheme,
		Host:   c.Request.Host,
		Path:   c.Request.URL.Path,
	}
	return u.String()
}

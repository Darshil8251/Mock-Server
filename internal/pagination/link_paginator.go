package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"mock-server/internal/config"
	"net/http"
	"net/url"
	"strconv"

	"mock-server/pkg/logger"

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
	var tmpLogger = logger.GetLogger()

	l := &linkPaginator{}

	responseObj, err := loadResponseObj(endpoint.ResponseObjFilePath)
	if err != nil {
		errInvalidResponse := fmt.Errorf("invalid response file path for endpoint: %s", endpoint.Path)
		tmpLogger.Warn(errInvalidResponse.Error(), err)
		return nil, errors.Join(errInvalidResponse, err)
	}

	l.responseObj = responseObj

	l.linkKey = defaultLinkKey

	l.paginationParameters = loadPaginationParameters(endpoint)

	if _, ok := endpoint.Pagination.Options["linkKey"].(string); ok {
		_, ok := responseObj[endpoint.Pagination.Options["linkKey"].(string)]
		if !ok {
			errInvalidLinkKey := fmt.Errorf("invalid link key for the endpoint: %v", endpoint.Path)
			tmpLogger.Warn(errInvalidLinkKey.Error(), err)
			return nil, errors.Join(errInvalidLinkKey, err)
		}
		l.linkKey = endpoint.Pagination.Options["linkKey"].(string)
	}

	// Validate the response field
	if endpoint.ResponseField != "" {
		_, ok := responseObj[endpoint.ResponseField].([]any)
		if !ok {
			errInvalidResponseField := fmt.Errorf("invalid response field for endpoint: %v", endpoint.Path)
			tmpLogger.Warn(errInvalidResponseField.Error(), err)
			return nil, errors.Join(errInvalidResponseField, err)
		}
		l.responseField = endpoint.ResponseField
		return l, nil
	}

	// If user not specified the response field, then find array field from the response object
	if endpoint.ResponseField == "" {
		for k, v := range responseObj {
			if _, ok := v.([]interface{}); ok {
				l.responseField = k
				return l, nil
			}
		}
	}

	errInvalidResponseField := fmt.Errorf("response field not present in response object for endpoint: %v", endpoint.Path)
	tmpLogger.Warn(errInvalidResponseField.Error(), err)
	return nil, errors.Join(errInvalidResponseField, err)
}

func (l *linkPaginator) Paginate(c *gin.Context) {
	var (
		pageSize = defaultPageSize
	)

	if l.paginationParameters.pageSentCount >= l.paginationParameters.totalPageCount {
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}

	if v := c.Query(l.paginationParameters.pageSizeKey); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			pageSize = p
		}
	}

	arr, ok := l.responseObj[l.responseField].([]any)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response field"})
		return
	}

	object := arr[0]
	APIResponseObject := make([]any, 0, pageSize)

	if l.paginationParameters.sentRecordsCount+pageSize > l.paginationParameters.totalRecordCount {
		pageSize = l.paginationParameters.totalRecordCount - l.paginationParameters.sentRecordsCount
	}

	for len(APIResponseObject) < pageSize {
		APIResponseObject = append(APIResponseObject, object)
	}

	l.paginationParameters.pageSentCount++
	l.paginationParameters.sentRecordsCount += pageSize

	l.responseObj[l.responseField] = APIResponseObject
	l.responseObj[l.linkKey] = generatePageLink(c)

	jsonResponse, err := json.Marshal(l.responseObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create response object"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResponse)

}

func generatePageLink(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	u := &url.URL{
		Scheme:   scheme,
		Host:     c.Request.Host,
		Path:     c.Request.URL.Path,
		RawQuery: c.Request.URL.RawQuery,
	}

	values, _ := url.ParseQuery(u.RawQuery)

	u.RawQuery = values.Encode()
	return u.String()
}

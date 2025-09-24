package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"mock-server/internal/config"
)

const (
	defaultPageKey          = "page"
	defaultPageSizeKey      = "pageSize"
	defaultOffsetKey        = "offset"
	defaultLimitKey         = "limit"
	defaultLinkKey          = "link"
	defaultPageSize         = 100
	defaultPageCount        = 2
	defaultTotalRecordCount = 200
)

// loadPaginationParameters loads the pagination parameters
func loadPaginationParameters(endpoint config.Endpoint) (p paginationParameters) {
	p = paginationParameters{}

	// Initialize with default values
	p.totalPageCount = defaultPageCount
	p.totalRecordCount = defaultTotalRecordCount
	p.pageSize = defaultPageSize
	p.pageKey = defaultPageKey
	p.pageSizeKey = defaultPageSizeKey
	p.pageSentCount = 0
	p.sentRecordsCount = 0

	if pageKey, ok := endpoint.Pagination.Options["pageKey"].(string); ok {
		p.pageKey = pageKey
	}

	if pageSizeKey, ok := endpoint.Pagination.Options["pageSizeKey"].(string); ok {
		p.pageSizeKey = pageSizeKey
	}

	if pageSize, ok := endpoint.Pagination.Options["pageSize"].(int); ok {
		p.pageSize = pageSize
	}

	if pageCount, ok := endpoint.Pagination.Options["totalPage"].(int); ok {
		p.totalPageCount = pageCount
	}

	if totalRecordCount, ok := endpoint.Pagination.Options["totalRecord"].(int); ok {
		p.totalRecordCount = totalRecordCount
	}

	p.pageParamsLocation = pageParameterLocation(endpoint.Pagination.Location)

	return p
}

// loadResponseObj loads the response object from the given file path
func loadResponseObj(path string) (responseObject map[string]any, err error) {
	if path == "" {
		return nil, errors.New("empty file path")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &responseObject); err != nil {
		return nil, err
	}

	return responseObject, nil
}

func findResponseFieldName(fieldName string, responseObject map[string]any) (string, error) {

	if fieldName == "" {
		for k, v := range responseObject {
			if _, ok := v.([]any); ok {
				fieldName = k
				break
			}
		}
		if fieldName == "" {
			return "", fmt.Errorf("response field doesn't exist in response object")
		}
		return fieldName, nil
	}

	_, ok := responseObject[fieldName].([]any)
	if !ok {
		return "", fmt.Errorf("invalid response field name")
	}

	return fieldName, nil

}

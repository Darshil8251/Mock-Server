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
	p.totalRecord = defaultTotalRecordCount
	p.pageSize = defaultPageSize
	p.sendPageCount = 0
	p.sendRecordsCount = 0

	fmt.Printf("Loading pagination parameters for endpoint: %s\n", endpoint.Pagination.Options)

	switch v := endpoint.Pagination.Options["pageSize"].(type) {
	case float64:
		p.pageSize = int(v)
	case int:
		p.pageSize = v
	}

	switch v := endpoint.Pagination.Options["totalRecord"].(type) {
	case float64:
		p.totalRecord = int(v)
	case int:
		p.totalRecord = v
	}

	// Calculate total pages using integer arithmetic
	p.totalPageCount = (p.totalRecord + p.pageSize - 1) / p.pageSize

	fmt.Println("Pagination Parameters Loaded:", p)

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

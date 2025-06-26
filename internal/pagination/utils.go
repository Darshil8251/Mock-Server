package pagination

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

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

	if pageCount, ok := endpoint.Pagination.Options["pageCount"].(int); ok {
		p.totalPageCount = pageCount
	}

	if totalRecordCount, ok := endpoint.Pagination.Options["totalRecordCount"].(int); ok {
		p.totalRecordCount = totalRecordCount
	}

	return p
}

// loadResponseObj loads the response object from the given file path
func loadResponseObj(path string) (map[string]interface{}, error) {
	if path == "" {
		return nil, errors.New("empty file path")
	}

	// Clean the path (removes ./ ../ etc.)
	cleanPath := filepath.Clean(path)

	responseObjFilePath := filepath.Join(".././", cleanPath)

	file, err := os.Open(responseObjFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(bytes, &responseObj); err != nil {
		return nil, err
	}

	return responseObj, nil
}

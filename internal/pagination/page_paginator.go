package pagination

import (
	"encoding/json"
	"fmt"
	"io"
	"mock-server/internal/config"
	"mock-server/pkg/logger"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func pagePaginator(endpoint config.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		tmpLogger := logger.GetLogger()
		// Extract config keys and defaults
		pageKey, _ := endpoint.Pagination.Options["pageKey"].(string)
		sizeKey, _ := endpoint.Pagination.Options["sizeKey"].(string)
		defaultPage, _ := endpoint.Pagination.Options["defaultPage"].(float64)
		defaultSize, _ := endpoint.Pagination.Options["defaultPageSize"].(float64)
		location := endpoint.Pagination.Location

		if pageKey == "" {
			pageKey = "page"
		}
		if sizeKey == "" {
			sizeKey = "size"
		}
		if defaultPage == 0 {
			defaultPage = 1
		}
		if defaultSize == 0 {
			defaultSize = 10
		}

		page := int(defaultPage)
		size := int(defaultSize)

		// Extract pagination params from the correct location
		switch location {
		case "body":
			var body map[string]interface{}
			_ = c.ShouldBindJSON(&body)
			if v, ok := body[pageKey].(float64); ok {
				page = int(v)
			} else if v, ok := body[pageKey].(int); ok {
				page = v
			}
			if v, ok := body[sizeKey].(float64); ok {
				size = int(v)
			} else if v, ok := body[sizeKey].(int); ok {
				size = v
			}
		case "header":
			if v := c.GetHeader(pageKey); v != "" {
				if p, err := strconv.Atoi(v); err == nil && p > 0 {
					page = p
				}
			}
			if v := c.GetHeader(sizeKey); v != "" {
				if s, err := strconv.Atoi(v); err == nil && s > 0 {
					size = s
				}
			}
		case "query":
			pageStr := c.DefaultQuery(pageKey, strconv.Itoa(int(defaultPage)))
			sizeStr := c.DefaultQuery(sizeKey, strconv.Itoa(int(defaultSize)))
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
			if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
				size = s
			}
		}

		if strings.Compare(endpoint.MockResponse, "") != 0 {
			endpoint.MockResponse = ".././" + endpoint.MockResponse
		}

		fmt.Println(endpoint.MockResponse)

		// Get the data to paginate
		file, err := os.Open(endpoint.MockResponse)
		if err != nil {
			tmpLogger.Error("failed to load response file", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load response file"})
			return
		}

		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response file"})
			return
		}

		var responseObj map[string]interface{}
		if err := json.Unmarshal(bytes, &responseObj); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response JSON"})
			return
		}

		// 2. Find the array field
		arrayField := endpoint.ArrayField
		if arrayField == "" {
			for k, v := range responseObj {
				if _, ok := v.([]interface{}); ok {
					arrayField = k
					break
				}
			}
		}
		if arrayField == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No array field found to update"})
			return
		}

		// 3. Get the array to paginate
		arr, _ := responseObj[arrayField].([]interface{})

		total := len(arr)
		start := (page - 1) * size
		end := start + size
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		pagedData := arr[start:end]

		c.JSON(http.StatusOK, gin.H{
			"items":       pagedData,
			"page":        page,
			"page_size":   size,
			"total_items": total,
			"total_pages": (total + size - 1) / size,
			"has_next":    end < total,
			"has_prev":    page > 1,
		})
	}
}

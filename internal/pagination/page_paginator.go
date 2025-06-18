package pagination

import (
	"mock-server/internal/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func pagePaginator(endpoint config.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		default: // "query" or fallback
			pageStr := c.DefaultQuery(pageKey, strconv.Itoa(int(defaultPage)))
			sizeStr := c.DefaultQuery(sizeKey, strconv.Itoa(int(defaultSize)))
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
			if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
				size = s
			}
		}

		// Get the data to paginate
		data, ok := endpoint.MockResponse.([]interface{})
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"items":       []interface{}{},
				"page":        page,
				"page_size":   size,
				"total_items": 0,
				"total_pages": 0,
				"has_next":    false,
				"has_prev":    false,
			})
			return
		}

		total := len(data)
		start := (page - 1) * size
		end := start + size
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		pagedData := data[start:end]

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

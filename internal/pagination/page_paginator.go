package pagination

import (
	"encoding/json"
	"fmt"
	"io"
	"mock-server/internal/config"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageKey  = "page"
	defaultPageSize = 100
)

func pagePaginator(endpoint config.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract config keys and defaults
		sizeKey, _ := endpoint.Pagination.Options["sizeKey"].(string)
		location := endpoint.Pagination.Location

		var pageSize = defaultPageSize

		fmt.Printf("sizeKey: %s\n", sizeKey)

		// Extract pagination params from the respective location
		switch PageParamsLocation(location) {
		case body:
			fmt.Printf("body: %v\n", c.Request.Body)
			var body map[string]interface{}
			err := c.ShouldBindJSON(&body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request body"})
				return
			}
			value, found := body[sizeKey]
			if found {
				pageSize = value.(int)
			}

		case header:
			if v := c.GetHeader(sizeKey); v != "" {
				if p, err := strconv.Atoi(v); err == nil && p > 0 {
					pageSize = p
				}
			}
		case query:
			size, err := strconv.Atoi(c.DefaultQuery(sizeKey, strconv.Itoa(defaultPageSize)))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sizeValue"})
				return
			}
			pageSize = size
		}

		if strings.Compare(endpoint.ResponseObjFilePath, "") != 0 {
			endpoint.ResponseObjFilePath = ".././" + endpoint.ResponseObjFilePath
		}

		// Get the data to paginate
		file, err := os.Open(endpoint.ResponseObjFilePath)
		if err != nil {
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
		arrayField := endpoint.ResponseField

		// If user not specified the response field, then find array field from the response object
		if strings.Compare(endpoint.ResponseField, "") == 0 {
			for k, v := range responseObj {
				if _, ok := v.([]interface{}); ok {
					arrayField = k
					break
				}
			}
		}
		if strings.Compare(arrayField, "") == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Response field not available"})
			return
		}

		// 3. Find the response object
		arr, ok := responseObj[arrayField].([]any)
		if !ok {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find response object"})
			return
		}

		object := arr[0]
		APIResponseObject := make([]any, 0, pageSize)

		for len(APIResponseObject) < pageSize {
			APIResponseObject = append(APIResponseObject, object)
		}

		responseObj[arrayField] = APIResponseObject

		jsonResponse, err := json.Marshal(responseObj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate the response"})
			return
		}

		c.Data(http.StatusOK, "application/json", jsonResponse)
	}
}

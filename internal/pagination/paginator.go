package pagination

import (
	"mock-server/internal/config"

	"github.com/gin-gonic/gin"
)

type PaginationType string

const (
	page   PaginationType = "page"
	token  PaginationType = "token"
	none   PaginationType = "none"
	link   PaginationType = "link"
	offset PaginationType = "offset"
)

func CreatePaginator(endpoint config.Endpoint) gin.HandlerFunc {

	switch PaginationType(endpoint.Pagination.Type) {

	case page:
		return pagePaginator(endpoint)
	default:
		return nil
	}
}

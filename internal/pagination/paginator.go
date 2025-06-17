package pagination

import (
	"fmt"
	
	"mock-server/internal/config"

	"github.com/gin-gonic/gin"
)

type paginationType string

type pageParamsLocation string

const (
	page   paginationType = "page"
	token  paginationType = "token"
	none   paginationType = "none"
	link   paginationType = "link"
	offset paginationType = "offset"

	body   pageParamsLocation = "body"
	query  pageParamsLocation = "query"
	header pageParamsLocation = "header"
)

type Paginator interface {
	Paginate(c *gin.Context)
}

func CreatePaginator(endpoint config.Endpoint) (Paginator, error) {

	switch paginationType(endpoint.Pagination.Type) {
	case page:
		p, err := createPagePaginator(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create page paginator for endpoint: %s", endpoint.Path)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unsupported pagination type: %s", endpoint.Pagination.Type)
	}
}

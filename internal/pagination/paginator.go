package pagination

import (
	"errors"
	"fmt"

	"mock-server/internal/config"

	"github.com/gin-gonic/gin"
)

type paginationType string

type pageParameterLocation string

const (
	page   paginationType = "page"
	token  paginationType = "token"
	none   paginationType = "none"
	link   paginationType = "link"
	offset paginationType = "offset"

	body   pageParameterLocation = "body"
	query  pageParameterLocation = "query"
	header pageParameterLocation = "header"
)

type Paginator interface {
	Paginate(c *gin.Context)
}

// paginationParameters use to keep tract of the pagination parameters
type paginationParameters struct {
	pageParamsLocation pageParameterLocation
	totalPageCount     int
	totalRecordCount   int
	pageKey            string
	pageSizeKey        string
	pageSentCount      int
	pageSize           int
	sentRecordsCount   int
}

func CreatePaginator(endpoint config.Endpoint) (Paginator, error) {

	switch paginationType(endpoint.Pagination.Type) {
	case page:
		p, err := createPagePaginator(endpoint)
		if err != nil {
			return nil, errors.Join(errors.New("failed to create page pagination"), err)
		}
		return p, nil
	case offset:
		p, err := createOffsetPaginator(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create offset paginator for endpoint: %s", endpoint.Path)
		}
		return p, nil
	case link:
		p, err := createLinkPaginator(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create link paginator for endpoint: %s", endpoint.Path)
		}
		return p, nil
	case token:
		p, err := createTokenPaginator(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create token paginator for endpoint : %s", endpoint.Path)
		}
		return p, nil
	default:
		p, err := createDefaultPaginator(endpoint)
		if err != nil {
			return nil, errors.Join(errors.New("failed to create default pagination"), err)
		}
		return p, nil
	}
}

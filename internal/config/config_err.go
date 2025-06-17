package config

import "errors"

var (
	ErrInvalidConfig         = errors.New("invalid configuration")
	ErrInvalidPath           = errors.New("invalid endpoint path")
	ErrInvalidMethod         = errors.New("invalid HTTP method for endpoint")
	ErrInvalidResponseFormat = errors.New("invalid response format for endpoint")
	ErrMissingHeaders        = errors.New("missing headers for endpoint")
	ErrMissingQueryParams    = errors.New("missing query parameters for endpoint")
	ErrMissingRequestBody    = errors.New("missing request body for endpoint")
	ErrMissingRateLimit      = errors.New("missing rate limit for endpoint")
	ErrMissingPagination     = errors.New("missing pagination configuration for endpoint")
	ErrMissingEndpoints      = errors.New("missing endpoints in configuration")
)

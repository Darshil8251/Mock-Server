package config

import "errors"

var (
	errInvalidConfig         = errors.New("invalid configuration")
	errInvalidPath           = errors.New("invalid endpoint path")
	errInvalidMethod         = errors.New("invalid HTTP method for endpoint")
	errInvalidResponseFormat = errors.New("invalid response format for endpoint")
	errMissingHeaders        = errors.New("missing headers for endpoint")
	errMissingQueryParams    = errors.New("missing query parameters for endpoint")
	errMissingRequestBody    = errors.New("missing request body for endpoint")
	errMissingRateLimit      = errors.New("missing rate limit for endpoint")
	errMissingPagination     = errors.New("missing pagination configuration for endpoint")
	errMissingEndpoints      = errors.New("missing endpoints in configuration")
)

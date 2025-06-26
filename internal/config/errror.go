package config

import "errors"

var (
	errInvalidConfig = errors.New("invalid configuration")
	errInvalidPath   = errors.New("invalid endpoint path")
	errInvalidMethod = errors.New("invalid HTTP method for endpoint")
)

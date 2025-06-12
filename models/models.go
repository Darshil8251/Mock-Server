package models

// Config represents the root configuration structure
type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Endpoint represents a single API endpoint configuration
type Endpoint struct {
	Path           string            `yaml:"path"`
	Method         string            `yaml:"method"`
	Headers        map[string]string `yaml:"headers"`
	QueryParams    map[string]string `yaml:"queryParams"`
	RequestBody    interface{}       `yaml:"requestBody"`
	Response       Response          `yaml:"response"`
	Pagination     Pagination        `yaml:"pagination"`
	ResponseParams ResponseParams    `yaml:"responseParams"`
}

// Response represents the expected response structure
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    interface{}       `yaml:"body"`
}

// Pagination represents pagination configuration
type Pagination struct {
	Enabled    bool   `yaml:"enabled"`
	PageSize   int    `yaml:"pageSize"`
	TotalItems int    `yaml:"totalItems"`
	PageParam  string `yaml:"pageParam"`
	SizeParam  string `yaml:"sizeParam"`
}

// ResponseParams defines which parameters should change in the response for each page
type ResponseParams struct {
	DynamicFields []string `yaml:"dynamicFields"`
	IncrementBy   int      `yaml:"incrementBy"`
}

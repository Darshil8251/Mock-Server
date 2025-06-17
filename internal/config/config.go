package config

import (
	"encoding/json"

	"mock-server/pkg/logger"
	"os"
)

type APIConfig struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path         string         `json:"path"`
	Method       string         `json:"method"`
	Headers      map[string]any `json:"headers"`
	QueryParams  map[string]any `json:"queryParams"`
	RequestBody  map[string]any `json:"requestBody"`
	RateLimit    int            `json:"rateLimit"`
	Pagination   Pagination     `json:"pagination"`
	MockResponse any            `json:"APIResponse"`
}

type Pagination struct {
	Type string `json:"type"`
}

func LoadConfig(configFilePath string) (*APIConfig, error) {
	apiConfig := &APIConfig{}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonParser := json.NewDecoder(jsonFile)
	if err := jsonParser.Decode(apiConfig); err != nil {
		return nil, err
	}

	err = validateConfig(apiConfig)
	if err != nil {
		logger.Get().ErrorW("Configuration validation failed", logger.Field("Err", err))
		return nil, err
	}

	return apiConfig, nil
}

func validateConfig(cfg *APIConfig) error {
	var tmpLogger = logger.Get()

	if cfg == nil {
		tmpLogger.WarnW("Provided configuration is nil", logger.Field("Err", ErrInvalidConfig))
		return nil
	}

	for _, endpoint := range cfg.Endpoints {
		if endpoint.Path == "" {
			tmpLogger.WarnW("Endpoint path cannot be empty", logger.Field("Err", ErrInvalidPath))
			return ErrInvalidPath
		}
		if endpoint.Method == "" {
			tmpLogger.WarnW("Endpoint method cannot be empty", logger.Field("Err", ErrInvalidMethod))
			return ErrInvalidMethod
		}
		if endpoint.MockResponse == nil {
			tmpLogger.WarnW("Mock response cannot be nil", logger.Field("Err", ErrInvalidResponseFormat))
			return ErrInvalidResponseFormat
		}
		if len(endpoint.Headers) == 0 {
			tmpLogger.WarnW("Headers cannot be empty", logger.Field("Err", ErrMissingHeaders))
			return ErrMissingHeaders
		}
		if len(endpoint.QueryParams) == 0 {
			tmpLogger.WarnW("Query parameters cannot be empty", logger.Field("Err", ErrMissingQueryParams))
			return ErrMissingQueryParams
		}
	}

	return nil
}

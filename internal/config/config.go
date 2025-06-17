package config

import (
	"encoding/json"
	"os"

	"mock-server/pkg/logger"
)

type APIConfig struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Path                string         `json:"path"`
	Method              string         `json:"method"`
	Headers             map[string]any `json:"headers"`
	QueryParams         map[string]any `json:"queryParams"`
	RequestBody         map[string]any `json:"requestBody"`
	RateLimit           int            `json:"rateLimit"`
	Pagination          pagination     `json:"pagination"`
	ResponseObjFilePath string         `json:"responseObjFilePath"`
	ResponseField       string         `json:"responseField,omitempty"`
}

type pagination struct {
	Type     string         `json:"type"`
	Location string         `json:"location"`
	Options  map[string]any `json:"options"`
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
		logger.GetLogger().Error("Configuration validation failed", errInvalidConfig)
		return nil, err
	}

	return apiConfig, nil
}

func validateConfig(cfg *APIConfig) error {
	var tmpLogger = logger.GetLogger()

	if cfg == nil {
		tmpLogger.Warn("Provided configuration is nil", errInvalidConfig)
		return nil
	}

	for _, endpoint := range cfg.Endpoints {
		if endpoint.Path == "" {
			tmpLogger.Warn("Endpoint path cannot be empty", errInvalidPath)
			return errInvalidPath
		}
		if endpoint.Method == "" {
			tmpLogger.Warn("Endpoint method cannot be empty", errInvalidMethod)
			return errInvalidMethod
		}
	}

	return nil
}

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
	Options  map[string]any `json:"options,omitempty"`
}

func LoadConfig() (*APIConfig, error) {
	var (
		mockLogger     = logger.GetLogger()
		configFilePath = os.Getenv("CONFIG_FILE_PATH")
	)

	apiConfig := &APIConfig{}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		mockLogger.Error("error opening config file", err)
		return nil, err
	}
	defer jsonFile.Close()

	jsonParser := json.NewDecoder(jsonFile)
	if err := jsonParser.Decode(apiConfig); err != nil {
		mockLogger.Error("error decoding config file", err)
		return nil, err
	}

	err = validateConfig(apiConfig)
	if err != nil {
		mockLogger.Error("invalid config", err)
		return nil, err
	}

	return apiConfig, nil
}

func validateConfig(cfg *APIConfig) error {
	var mockLogger = logger.GetLogger()

	if cfg == nil {
		mockLogger.Warn("provided config is nil", errInvalidConfig)
		return nil
	}

	for _, endpoint := range cfg.Endpoints {
		if endpoint.Path == "" {
			mockLogger.Warn("invalid endpoint path", errInvalidPath)
			return errInvalidPath
		}
		if endpoint.Method == "" {
			mockLogger.Warn("invalid endpoint method", errInvalidMethod)
			return errInvalidMethod
		}
	}

	return nil
}

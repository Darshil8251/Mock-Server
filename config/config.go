package config

import (
	"fmt"
	"os"

	"mock-server/models"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filePath string) (*models.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ValidateConfig validates the configuration
func ValidateConfig(config *models.Config) error {
	for _, endpoint := range config.Endpoints {
		if endpoint.Path == "" {
			return fmt.Errorf("endpoint path cannot be empty")
		}
		if endpoint.Method == "" {
			return fmt.Errorf("endpoint method cannot be empty")
		}
	}
	return nil
}

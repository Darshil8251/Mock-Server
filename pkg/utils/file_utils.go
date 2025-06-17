package utils

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Helper function to handle path resolution
func ValidateAndResolvePath(rawPath string) (string, error) {
	if rawPath == "" {
		return "", errors.New("empty file path")
	}

	// Clean the path (removes ./ ../ etc.)
	cleanPath := filepath.Clean(rawPath)

	return filepath.Join(".././", cleanPath), nil
}

// Helper function to read and parse JSON
func ReadAndParseJSONFile(path string) (map[string]interface{}, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(bytes, &responseObj); err != nil {
		return nil, err
	}

	return responseObj, nil
}

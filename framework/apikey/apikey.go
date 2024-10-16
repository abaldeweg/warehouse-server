package apikey

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

// APIKey represents a single API key.
type APIKey struct {
	Key         string   `json:"key"`
	Permissions []string `json:"permissions"`
}

// keysFilePath is the path to the JSON file storing API keys.
const keysFilePath = "data/auth/api_keys.json"

// loadAPIKeys loads API keys from the JSON file.
func loadAPIKeys() ([]APIKey, error) {
	data, err := os.ReadFile(keysFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []APIKey{}, createAPIKeysFile(nil)
		}

		return nil, fmt.Errorf("error reading API keys file: %w", err)
	}

	var keys []APIKey
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("error unmarshalling API keys: %w", err)
	}

	return keys, nil
}

// createAPIKeysFile saves API keys to the JSON file.
// If no keys are provided, it generates a new UUID as a default key.
func createAPIKeysFile(keys []APIKey) error {
	if len(keys) == 0 {
		newKey := APIKey{Key: uuid.New().String(), Permissions: []string{}}
		keys = append(keys, newKey)
	}

	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling API keys: %w", err)
	}

	return os.WriteFile(keysFilePath, data, 0644)
}

// IsValidAPIKey checks if the given key is valid.
func IsValidAPIKey(key string) bool {
	keys, err := loadAPIKeys()
	if err != nil {
		fmt.Println("Error loading API keys:", err)

		return false
	}

	for _, k := range keys {
		if k.Key == key {
			return true
		}
	}

	return false
}

// HasPermission checks if the given API key has the specified permission.
func HasPermission(key string, permission string) bool {
	keys, err := loadAPIKeys()
	if err != nil {
		fmt.Println("Error loading API keys:", err)
		return false
	}

	for _, k := range keys {
		if k.Key == key {
			for _, p := range k.Permissions {
				if p == permission {
					return true
				}
			}
		}
	}

	return false
}

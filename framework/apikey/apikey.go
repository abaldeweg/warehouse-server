package apikey

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// APIKeys represents a collection of API keys.
type APIKeys struct {
	Keys []APIKey `json:"keys"`
}

// APIKey represents a single API key with associated permissions.
type APIKey struct {
	Key         string   `json:"key"`
	Permissions []string `json:"permissions"`
}

// NewAPIKeys creates a new APIKeys instance from JSON data.
func NewAPIKeys(data []byte) (*APIKeys, error) {
	var keys APIKeys
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("error unmarshalling API keys: %w", err)
	}

	return &keys, nil
}

// DefaultAPIKeysFile generates the default API keys file content.
func DefaultAPIKeysFile() ([]byte, error) {
	keys := APIKeys{Keys: []APIKey{{Key: uuid.New().String(), Permissions: []string{}}}}

	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling API keys: %w", err)
	}

	return data, nil
}

// IsValidAPIKey checks if the given key exists in the APIKeys collection.
func (keys *APIKeys) IsValidAPIKey(key string) bool {
	for _, k := range keys.Keys {
		if k.Key == key {
			return true
		}
	}

	return false
}

// HasPermission checks if the given API key has the specified permission.
func (keys *APIKeys) HasPermission(key string, permission string) bool {
	for _, k := range keys.Keys {
		if k.Key == key {
			for _, p := range k.Permissions {
				if p == permission {
					return true
				}
			}
			return false
		}
	}

	return false
}

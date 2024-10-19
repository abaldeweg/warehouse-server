package apikey

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIKeysWithValidJSON(t *testing.T) {
	jsonData := []byte(`{"keys": [{"key": "test-key", "permissions": ["read"]}]}`)
	keys, err := NewAPIKeys(jsonData)

	assert.NoError(t, err)
	assert.NotNil(t, keys)
	assert.Equal(t, 1, len(keys.Keys))
	assert.Equal(t, "test-key", keys.Keys[0].Key)
	assert.Equal(t, []string{"read"}, keys.Keys[0].Permissions)
}

func TestNewAPIKeysWithInvalidJSON(t *testing.T) {
	jsonData := []byte(`invalid json`)
	keys, err := NewAPIKeys(jsonData)

	assert.Error(t, err)
	assert.Nil(t, keys)
}

func TestDefaultAPIKeysFile(t *testing.T) {
	data, err := DefaultAPIKeysFile()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var keys APIKeys
	err = json.Unmarshal(data, &keys)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys.Keys))
	assert.NotEmpty(t, keys.Keys[0].Key)
	assert.Empty(t, keys.Keys[0].Permissions)
}

func TestIsValidAPIKey(t *testing.T) {
	keys := &APIKeys{
		Keys: []APIKey{
			{Key: "valid-key", Permissions: []string{}},
		},
	}

	assert.True(t, keys.IsValidAPIKey("valid-key"))
	assert.False(t, keys.IsValidAPIKey("invalid-key"))
}

func TestHasPermission(t *testing.T) {
	keys := &APIKeys{
		Keys: []APIKey{
			{Key: "key-1", Permissions: []string{"read", "write"}},
			{Key: "key-2", Permissions: []string{"read"}},
		},
	}

	assert.True(t, keys.HasPermission("key-1", "read"))
	assert.True(t, keys.HasPermission("key-1", "write"))
	assert.True(t, keys.HasPermission("key-2", "read"))

	assert.False(t, keys.HasPermission("key-1", "delete"))
	assert.False(t, keys.HasPermission("key-2", "write"))
	assert.False(t, keys.HasPermission("invalid-key", "read"))
}

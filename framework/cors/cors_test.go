package cors

import (
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestDefaultConfig tests the DefaultConfig function.
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.Contains(t, config.AllowHeaders, "Authorization")
	assert.Equal(t, []string{"*"}, config.AllowOrigins)

	viper.Set("CORS_ALLOW_ORIGIN", "http://localhost,http://test.localhost")
	config = DefaultConfig()
	assert.Equal(t, []string{"http://localhost", "http://test.localhost"}, config.AllowOrigins)
}

// TestSetDefaultCorsHeaders tests the SetDefaultCorsHeaders function.
func TestSetDefaultCorsHeaders(t *testing.T) {
	handler := SetDefaultCorsHeaders()
	assert.NotNil(t, handler)
}

// TestSetCorsHeaders tests the SetCorsHeaders function.
func TestSetCorsHeaders(t *testing.T) {
	customConfig := cors.Config{
		AllowOrigins: []string{"http://localhost"},
	}
	handler := SetCorsHeaders(customConfig)
	assert.NotNil(t, handler)
}

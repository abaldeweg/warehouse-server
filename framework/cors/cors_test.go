package cors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCorsConfig(t *testing.T) {
	corsConfig := NewCors()

	assert.Equal(t, []string{"http://127.0.0.1"}, corsConfig.AllowOrigins)
	assert.Equal(t, []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}, corsConfig.AllowMethods)
	assert.Equal(t, []string{"Origin", "Authorization", "Content-Type"}, corsConfig.AllowHeaders)
	assert.Equal(t, []string{"Content-Length"}, corsConfig.ExposeHeaders)
}

func TestNewCorsConfigEdit(t *testing.T) {
	corsConfig := NewCors()

	corsConfig.AllowOrigins = []string{"http://test.localhost"}
	corsConfig.AllowMethods = []string{"POST", "GET", "OPTIONS", "PUT"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	corsConfig.ExposeHeaders = []string{"Content-Type"}

	assert.Equal(t, []string{"http://test.localhost"}, corsConfig.AllowOrigins)
	assert.Equal(t, []string{"POST", "GET", "OPTIONS", "PUT"}, corsConfig.AllowMethods)
	assert.Equal(t, []string{"Authorization", "Content-Type"}, corsConfig.AllowHeaders)
	assert.Equal(t, []string{"Content-Type"}, corsConfig.ExposeHeaders)
}

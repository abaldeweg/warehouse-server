package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadAppConfig(t *testing.T) {
	t.Run("Config File Not Found", func(t *testing.T) {
		LoadAppConfig(WithName("config"), WithFormat("json"), WithPaths("."))
		assert.Equal(t, "", viper.ConfigFileUsed())
	})

	t.Run("Valid Config File", func(t *testing.T) {
		LoadAppConfig(WithName("config"), WithFormat("json"), WithPaths("testdata"))
		expectedPath, _ := filepath.Abs("testdata/config.json")
		assert.Equal(t, expectedPath, viper.ConfigFileUsed())
	})
}

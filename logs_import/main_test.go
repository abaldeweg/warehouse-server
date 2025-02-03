package main

import (
	"testing"

	"github.com/spf13/viper"
)

func TestIsBlocked(t *testing.T) {
	viper.Set("blocklist", []string{"/blocked", "/spam"})

	tests := []struct {
		path     string
		expected bool
	}{
		{"/blocked", true},
		{"/spam", true},
		{"/", false},
	}

	for _, test := range tests {
		result := isBlocked(test.path)
		if result != test.expected {
			t.Errorf("isBlocked(%s) = %v; want %v", test.path, result, test.expected)
		}
	}
}

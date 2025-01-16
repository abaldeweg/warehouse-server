package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Option is a function that sets a configuration option.
type Option func()

// LoadAppConfig loads the application configuration.
func LoadAppConfig(options ...Option) {
	for _, option := range options {
		option()
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("The Config file was not found, using defaults.")
		} else {
			log.Fatalf("Error loading config file: %s", err)
		}
	}
}

// WithName sets the name of the configuration file.
func WithName(name string) Option {
	return func() {
		viper.SetConfigName(name)
	}
}

// WithFormat sets the format of the configuration file.
func WithFormat(format string) Option {
	return func() {
		viper.SetConfigType(format)
	}
}

// WithPaths sets the paths to search for the configuration file.
func WithPaths(paths ...string) Option {
	return func() {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}
}

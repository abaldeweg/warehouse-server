package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

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

type Option func()

func WithName(name string) Option {
	return func() {
		viper.SetConfigName(name)
	}
}

func WithFormat(format string) Option {
	return func() {
		viper.SetConfigType(format)
	}
}

func WithPaths(paths ...string) Option {
	return func() {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}
}

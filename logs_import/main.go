package main

import (
	"time"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/logs_import/cmd"
	"github.com/spf13/viper"
)

func main() {
	config.LoadAppConfig()

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")

	go cmd.Execute()

	for {
		time.Sleep(1 * time.Second)
	}
}

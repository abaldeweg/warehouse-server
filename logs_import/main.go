package main

import (
	"time"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/logs_import/importer"
	"github.com/spf13/viper"
)

func main() {
	config.LoadAppConfig(config.WithName("config"), config.WithFormat("json"), config.WithPaths("./data/config"))

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("blocklist", []string{})

	go importer.Import()

	for {
		time.Sleep(1 * time.Second)
	}
}

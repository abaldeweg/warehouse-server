package main

import (
	"time"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/logs_import/importer"
	"github.com/spf13/viper"

	"fmt"
)

func main() {
	config.LoadAppConfig(config.WithName("config"), config.WithFormat("json"), config.WithPaths("./data/config"))

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("blocklist", []string{})

	fmt.Println("Blocklist:", viper.GetStringSlice("blocklist"))

	go importer.ImportLogs()

	for {
		time.Sleep(1 * time.Second)
	}
}

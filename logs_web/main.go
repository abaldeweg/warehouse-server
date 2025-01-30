package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/logs_web/router"
	"github.com/spf13/viper"
)

func main() {
	config.LoadAppConfig()

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")

	r := router.Routes()
	log.Fatal(r.Run(":8080"))
}

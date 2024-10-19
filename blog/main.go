package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/blog/router"
	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/framework/cors"
	"github.com/spf13/viper"
)

func main() {
	config.LoadAppConfig()

    viper.SetDefault("CORS_ALLOW_ORIGIN", "*")

	r := router.Routes()
	r.Use(cors.SetDefaultCorsHeaders())

	log.Fatal(r.Run(":8080"))
}

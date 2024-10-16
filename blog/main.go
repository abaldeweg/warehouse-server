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

	corsConfig := cors.NewCors()
    corsConfig.AllowOrigins = []string{viper.GetString("CORS_ALLOW_ORIGIN")}
	corsConfig.SetCorsHeaders()

	r := router.Routes()
	r.Use(corsConfig.SetCorsHeaders())

	log.Fatal(r.Run(":8080"))
}

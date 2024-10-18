package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/framework/cors"
	"github.com/abaldeweg/warehouse-server/gateway/router"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	godotenv.Load()

	config.LoadAppConfig()

	viper.SetDefault("CORS_ALLOW_ORIGIN", "*")

	corsConfig := cors.NewCors()
	corsConfig.Config.AllowOrigins = []string{viper.GetString("CORS_ALLOW_ORIGIN")}
	corsConfig.SetCorsHeaders()

	r := router.Routes()
	r.Use(corsConfig.SetCorsHeaders())

	log.Fatal(r.Run(":8080"))
}

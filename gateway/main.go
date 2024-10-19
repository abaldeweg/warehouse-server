package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/abaldeweg/warehouse-server/framework/cors"
	"github.com/abaldeweg/warehouse-server/gateway/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.LoadAppConfig()

	r := router.Routes()
	r.Use(cors.SetDefaultCorsHeaders())

	log.Fatal(r.Run(":8080"))
}

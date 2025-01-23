package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/logs_web/router"
)

func main() {
	r := router.Routes()
	log.Fatal(r.Run(":8080"))
}

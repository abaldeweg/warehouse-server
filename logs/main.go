package main

import (
	"log"

	"github.com/abaldeweg/warehouse-server/logs/router"
)

func main() {
	r := router.Routes()
	log.Fatal(r.Run(":8080"))
}

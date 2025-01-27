package main

import (
	"time"

	"github.com/abaldeweg/warehouse-server/logs_import/cmd"
)

func main() {
	go cmd.Execute()

	for {
		time.Sleep(1 * time.Second)
	}
}

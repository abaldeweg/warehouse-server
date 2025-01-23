package main

import (
    "github.com/abaldeweg/warehouse-server/logs_import/cmd"
    "time"
)

func main() {
    go cmd.Execute()

    for {
        time.Sleep(1 * time.Second)
    }
}

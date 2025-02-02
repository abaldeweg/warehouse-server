package main

import (
	"time"

	"github.com/abaldeweg/warehouse-server/framework/config"
	"github.com/spf13/viper"

	"fmt"
	"log"
	"os"

	"github.com/abaldeweg/warehouse-server/logs_import/db"
	"github.com/abaldeweg/warehouse-server/logs_import/parser"
)

func main() {
  config.LoadAppConfig(config.WithName("config"), config.WithFormat("json"), config.WithPaths("./data/config"))

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("blocklist", []string{})

	go importLogs()

	for {
		time.Sleep(1 * time.Second)
	}
}

func importLogs() {
	entries, err := parser.ReadLogEntries()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := db.NewDBHandler()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	for _, entry := range entries {
		if !isBlocked(entry.RequestPath) {
			if err := db.Add(entry); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}
	}

	fmt.Println("\033[32mLogs successfully imported!\033[0m")
	os.Exit(0)
}

func isBlocked(host string) bool {
	blocklist := viper.GetStringSlice("blocklist")
	for _, blocked := range blocklist {
		if host == blocked {
			return true
		}
	}
	return false
}

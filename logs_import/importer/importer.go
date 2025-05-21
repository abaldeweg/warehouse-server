package importer

import (
	"fmt"
	"log"
	"os"

	"github.com/abaldeweg/warehouse-server/logs_import/db"
	"github.com/abaldeweg/warehouse-server/logs_import/parser"
	"github.com/spf13/viper"
)

func Import() {
	if err := importLogs(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("\033[32mLogs successfully imported!\033[0m")
}

func importLogs() error {
	entries, err := parser.ReadLogEntries()
	if err != nil {
		return err
	}

	db, err := db.NewDBHandler()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, entry := range entries {
		if !isBlocked(entry.RequestPath) {
			if err := db.Add(entry); err != nil {
				return err
			}
		}
	}
	return nil
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

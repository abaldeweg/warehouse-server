package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/abaldeweg/warehouse-server/logs_import/db"
	"github.com/abaldeweg/warehouse-server/logs_import/parser"
)

// ImportLogsCmd reads logs from the log file and imports them into the database.
var ImportLogsCmd = &cobra.Command{
	Use:   "import",
	Short: "Read logs from the log file and import them into the database",
	Run: func(cmd *cobra.Command, args []string) {
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
			if err := db.Add(entry); err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}

		fmt.Println("\033[32mLogs successfully imported!\033[0m")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(ImportLogsCmd)
}

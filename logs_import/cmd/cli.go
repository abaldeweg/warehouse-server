package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/abaldeweg/warehouse-server/logs/db"
	"github.com/abaldeweg/warehouse-server/logs/parser"
)

// ImportLogsCmd reads logs from the log file and imports them into the database.
var ImportLogsCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Read logs from the log file and import them into the database",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var fileName string
		if len(args) > 0 {
			fileName = args[0]
		}

		entries, err := parser.ReadLogEntries(fileName)
		if err != nil {
			log.Fatal(err)
		}

		h, err := db.NewDBHandler()
		if err != nil {
			log.Fatal(err)
		}
		defer h.Close()

		for _, entry := range entries {
			date, _ := strconv.Atoi(time.Time(entry.Time).Format("20060102"))
			exists, err := h.Exists(date, entry)
			if err != nil {
				log.Fatal(err)
			}
			if !exists {
				if err := h.Write(date, entry); err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Println("\033[32mLogs successfully imported!\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(ImportLogsCmd)
}

package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/abaldeweg/warehouse-server/logs/db"
	"github.com/abaldeweg/warehouse-server/logs/parser"
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

		h, err := db.NewDBHandler()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		defer h.Close()

		for _, entry := range entries {
			date, _ := strconv.Atoi(time.Time(entry.Time).Format("20060102"))
			fmt.Println(date)
			exists, err := h.Exists(date, entry)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			if !exists {
				if err := h.Write(date, entry); err != nil {
					log.Println(err)
					os.Exit(1)
				}
			}
		}

		fmt.Println("\033[32mLogs successfully imported!\033[0m")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(ImportLogsCmd)
}

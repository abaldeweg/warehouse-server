package cmd

import (
	"github.com/abaldeweg/warehouse-server/logs/controller"
	"github.com/spf13/cobra"
)

// ImportLogsCmd reads logs from the log file and imports them into the database.
var ImportLogsCmd = &cobra.Command{
	Use:   "import",
	Short: "Read logs from the log file and import them into the database",
	Run: func(cmd *cobra.Command, args []string) {
		controller.ImportLogs()
	},
}

func init() {
	rootCmd.AddCommand(ImportLogsCmd)
}

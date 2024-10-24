package cmd

import (
	"fmt"
	"log"

	"github.com/abaldeweg/warehouse-server/warehousecli/command"
	"github.com/spf13/cobra"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot the OS.",
	Long:  `Reboot the OS. The connection will be terminated.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reboot system")
		fmt.Println("The connection will be terminated")

		_, err := command.Command([]string{"reboot"})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rebootCmd)
}

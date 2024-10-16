package cmd

import (
	"fmt"
	"log"

	"github.com/abaldeweg/warehouse-server/admincli/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Reload container images",
	Long:  `Restart all containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping all containers...")

		stopContainers, err := command.Command([]string{"/usr/bin/docker", "compose", "--project-directory", viper.GetString("project_dir"), "down"})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(stopContainers))
		fmt.Println("Done")

		fmt.Println("Starting all containers...")

		startContainers, err := command.Command([]string{"/usr/bin/docker", "compose", "--project-directory", viper.GetString("project_dir"), "up", "-d"})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(startContainers))
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

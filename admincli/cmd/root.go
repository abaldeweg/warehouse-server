package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "admincli",
	Short: "Maintenance tools",
	Long:  `The app gives you simple access to maintenance tools.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetDefault("project_dir", ".")

	viper.SetConfigName("admincli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/admincli/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

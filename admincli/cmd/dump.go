package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump database",
	Long:  `Dumps the database to the file system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dumping database")

		// command
		c := exec.Command("docker", "exec", viper.GetString("database"), "sh", "-c", fmt.Sprintf("exec mysqldump %s -uroot -p\"%s\"", "$MYSQL_DATABASE", "$MYSQL_ROOT_PASSWORD"))

		// create directory
		dirPath := fmt.Sprintf("%s/dump/", viper.GetString("project_dir"))

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		// create dump
		weekday := strings.ToLower(time.Now().Weekday().String())

		file, err := os.Create(fmt.Sprintf("%sdump_%s.sql", dirPath, weekday))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// run command
		c.Stdout = file
		err = c.Run()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Done")
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}

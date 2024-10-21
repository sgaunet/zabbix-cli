package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "development"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version of zabbix-cli",
	Long:  `print version of zabbix-cli`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

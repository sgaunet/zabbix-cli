package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var example string = `
zabbix_endpoint: http://zabbix.mydomain.com/api_jsonrpc.php
zabbix_user: admin
zabbix_password: *****
`

// PrintConfigCmd represents the config command
var PrintConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "describes how to configure zabbix-cli",
	Long:  `describes how to configure zabbix-cli`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Default configuration file: $HOME/.config/zabbix-cli/default.yaml")
		fmt.Println("Create $HOME/.config/zabbix-cli and add configuration files in it.")
		fmt.Println("")
		fmt.Println("Below is an example of configuration file:")
		fmt.Println(example)
	},
}

func init() {
	rootCmd.AddCommand(PrintConfigCmd)
}

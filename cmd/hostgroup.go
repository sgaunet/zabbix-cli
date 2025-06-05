package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// HostgroupCmd represents the get subcommand (get problems, get hosts ...)
var HostgroupCmd = &cobra.Command{
	Use:   "hostgroup",
	Short: "get hostgroups",
	Long:  `get hostgroups`,
	Run: func(cmd *cobra.Command, _ []string) {
		// print help
		cmd.Help() //nolint:errcheck
		os.Exit(1)
	},
}

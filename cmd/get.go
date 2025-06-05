package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// GetCmd represents the get subcommand (get problems, get hosts ...)
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run: func(cmd *cobra.Command, _ []string) {
		// print help
		cmd.Help() //nolint:errcheck
		os.Exit(1)
	},
}

func init() {
	GetCmd.AddCommand(GetHostGroupCmd) // GetHostGroupCmd is defined in get-hostgroup.go
	GetCmd.AddCommand(GetMaintenanceCmd) // GetMaintenanceCmd is defined in get-maintenance.go
}

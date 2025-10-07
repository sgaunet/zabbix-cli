package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// DashboardCmd represents the dashboard subcommand
var DashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "dashboard operations",
	Long:  `dashboard operations`,
	Run: func(cmd *cobra.Command, _ []string) {
		// print help
		cmd.Help() //nolint:errcheck
		os.Exit(1)
	},
}

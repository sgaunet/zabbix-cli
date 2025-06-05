package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// MaintenanceCmd represents the maintenance subcommand
var MaintenanceCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Create and manage maintenance periods",
	Long:  `Create and manage maintenance periods for hosts and host groups.`,
	Run: func(cmd *cobra.Command, _ []string) {
		// print help
		cmd.Help() //nolint:errcheck
		os.Exit(1)
	},
}

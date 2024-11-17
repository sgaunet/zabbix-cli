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
		cmd.Help()
		os.Exit(1)
	},
}

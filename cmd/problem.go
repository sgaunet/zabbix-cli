package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// ProblemCmd represents the get subcommand (get problems, get hosts ...)
var ProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "get problems",
	Long:  `get problems`,
	Run: func(cmd *cobra.Command, _ []string) {
		// print help
		cmd.Help() //nolint:errcheck
		os.Exit(1)
	},
}

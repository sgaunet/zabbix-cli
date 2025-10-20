package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var DashboardExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export dashboard configuration",
	Long:  `Export a dashboard configuration in YAML, JSON, or XML format. You can specify the dashboard by name or ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Validate that at least one of name or ID is provided
		if dashboardName == "" && dashboardID == "" {
			return fmt.Errorf("either --name or --id must be specified")
		}

		// Validate that both are not provided
		if dashboardName != "" && dashboardID != "" {
			return fmt.Errorf("cannot specify both --name and --id")
		}

		// Get format option
		formatOpt, err := GetFormatOption(dashboardExportFormat)
		if err != nil {
			return fmt.Errorf("invalid format: %w", err)
		}

		// initConfig() should be called by cobra.OnInitialize on the RootCmd.
		if conf == nil {
			return fmt.Errorf("configuration not initialized; initConfig did not run or failed")
		}

		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		if err := z.Login(ctx); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		defer func() {
			if err := z.Logout(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "logout failed: %v\n", err)
			}
		}()

		// Find dashboard by name or ID
		var dashboardIDToExport string
		if dashboardName != "" {
			// Search by name
			filter := map[string]interface{}{
				"name": dashboardName,
			}
			req := zabbix.NewDashboardGetRequest(
				zabbix.WithDashboardGetAuth(z.Auth()),
				zabbix.WithDashboardGetID(1),
				zabbix.WithDashboardGetOutput("extend"),
				zabbix.WithDashboardGetFilter(filter),
			)

			response, err := z.DashboardGet(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to get dashboard by name: %w", err)
			}

			if len(response.Result) == 0 {
				return fmt.Errorf("dashboard with name '%s' not found", dashboardName)
			}

			if len(response.Result) > 1 {
				return fmt.Errorf("multiple dashboards found with name '%s', please use --id instead", dashboardName)
			}

			dashboardIDToExport = response.Result[0].DashboardID
		} else {
			// Use provided ID
			dashboardIDToExport = dashboardID
		}

		// Export the dashboard
		res, err := z.Export(ctx,
			formatOpt,
			zabbix.ExportRequestOptionDashboardsID([]string{dashboardIDToExport}),
		)
		if err != nil {
			return fmt.Errorf("failed to export dashboard: %w", err)
		}

		// Output to file or stdout
		if dashboardFile != "" {
			err = os.WriteFile(dashboardFile, []byte(res), 0644)
			if err != nil {
				return fmt.Errorf("failed to write dashboard to file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Dashboard exported to %s\n", dashboardFile)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), res)
		}

		return nil
	},
}

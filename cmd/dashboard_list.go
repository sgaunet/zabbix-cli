package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var dashboardOutputFormat string

var DashboardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Zabbix dashboards",
	Long:  `List Zabbix dashboards. By default, displays results in a table. Use --output json for raw JSON.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// initConfig() should be called by cobra.OnInitialize on the RootCmd.
		// If conf is nil here, it means initConfig was not called or failed.
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

		// Prepare request to get all dashboards
		dashboardRequest := zabbix.NewGetAllDashboardsRequest(
			zabbix.WithDashboardGetAuth(z.Auth()),
			zabbix.WithDashboardGetID(1), // Simplified ID for CLI context
		)

		response, err := z.DashboardGet(ctx, dashboardRequest)
		if err != nil {
			return fmt.Errorf("failed to get dashboards: %w", err)
		}

		if dashboardOutputFormat == "json" {
			jsonOutput, err := json.MarshalIndent(response.Result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal dashboards to JSON: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(jsonOutput))
		} else {
			// Table output
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tPRIVATE\tDISPLAY_PERIOD")
			for _, dashboard := range response.Result {
				privateStr := "No"
				if dashboard.Private.Bool() {
					privateStr = "Yes"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					dashboard.DashboardID,
					dashboard.Name,
					privateStr,
					dashboard.DisplayPeriod)
			}
			w.Flush()
		}

		return nil
	},
}

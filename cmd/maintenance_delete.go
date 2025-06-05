package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

// MaintenanceDeleteCmd represents the delete maintenance subcommand
var MaintenanceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete maintenance periods",
	Long:  `Delete maintenance periods from your Zabbix installation.`,
	Run: func(cmd *cobra.Command, _ []string) {
		// Print help if no subcommands are provided
		cmd.Help() //nolint:errcheck
	},
}

// MaintenanceDeleteAllCmd represents the command to delete all maintenance periods
var MaintenanceDeleteAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Delete all maintenance periods",
	Long:  `Delete all maintenance periods from your Zabbix installation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		ctx := context.Background()

		err = initConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		err = z.Login(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}
		defer z.Logout(ctx) //nolint:errcheck

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// First, get all maintenance periods
		maintenanceRequest := zabbix.NewMaintenanceGetRequest(zabbix.WithMaintenanceGetAuthToken(z.Auth()))
		response, err := z.MaintenanceGet(ctx, maintenanceRequest)
		if err != nil {
			return fmt.Errorf("failed to get maintenance periods: %w", err)
		}

		if len(response.Result) == 0 {
			fmt.Println("No maintenance periods found to delete")
			return nil
		}

		// Extract maintenance IDs for deletion
		maintenanceIDs := make([]string, 0, len(response.Result))
		for _, maintenance := range response.Result {
			maintenanceIDs = append(maintenanceIDs, maintenance.MaintenanceID)
		}

		// Delete all maintenance periods
		deleteRequest := zabbix.NewMaintenanceDeleteRequest(
			zabbix.WithMaintenanceDeleteIDs(maintenanceIDs),
			zabbix.WithMaintenanceDeleteAuthToken(z.Auth()),
		)

		deleteResponse, err := z.MaintenanceDelete(ctx, deleteRequest)
		if err != nil {
			return fmt.Errorf("failed to delete maintenance periods: %w", err)
		}

		fmt.Printf("Successfully deleted %d maintenance periods\n", len(deleteResponse.Result.MaintenanceIDs))
		return nil
	},
}

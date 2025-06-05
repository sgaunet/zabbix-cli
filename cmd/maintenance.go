package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var (
	maintenanceName        string
	maintenanceDescription string
	maintenanceDuration    int
	maintenanceType        int
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

// MaintenanceCreateAllCmd represents the create-all subcommand that creates maintenance for all host groups
var MaintenanceCreateAllCmd = &cobra.Command{
	Use:   "create-all",
	Short: "Create a maintenance period for all host groups",
	Long:  `Create a one-time maintenance period for all host groups, active from now for the specified duration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Validate configuration
		if conf == nil {
			return fmt.Errorf("configuration not initialized; initConfig did not run or failed")
		}

		// Initialize Zabbix client
		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		if err := z.Login(ctx); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		defer func() {
			if err := z.Logout(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "logout failed: %v\n", err)
			}
		}()

		// Get all host groups
		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth(z.Auth()),
			zabbix.WithHostGroupGetID(1),
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}

		if len(response.Result) == 0 {
			return fmt.Errorf("no host groups found")
		}

		// Extract host group IDs
		groupIDs := make([]string, 0, len(response.Result))
		for _, hg := range response.Result {
			groupIDs = append(groupIDs, hg.GroupID)
		}

		// Calculate maintenance period
		now := time.Now()
		activeTill := now.Add(time.Duration(maintenanceDuration) * time.Hour).Unix()

		// Create a one-time maintenance period
		timePeriod := []zabbix.TimePeriod{
			{
				TimePeriodType: zabbix.TimePeriodTypeOneTime,
				StartTime:      0,                          // Start at the beginning of active_since
				Period:         maintenanceDuration * 3600, // Duration in seconds
			},
		}

		// Prepare maintenance creation request
		maintenanceRequest := zabbix.NewMaintenanceCreateRequest(
			zabbix.WithMaintenanceName(maintenanceName),
			zabbix.WithMaintenanceDescription(maintenanceDescription),
			zabbix.WithMaintenanceActiveTill(activeTill),
			zabbix.WithMaintenanceType(zabbix.MaintenanceType(maintenanceType)),
			zabbix.WithMaintenanceTimePeriods(timePeriod),
			zabbix.WithMaintenanceGroupIDs(groupIDs),
			zabbix.WithMaintenanceAuthToken(z.Auth()),
			zabbix.WithMaintenanceRequestID(1),
		)

		// Create maintenance
		maintenanceResponse, err := z.MaintenanceCreate(ctx, maintenanceRequest)
		if err != nil {
			return fmt.Errorf("failed to create maintenance: %w", err)
		}

		// Display result
		fmt.Fprintf(cmd.OutOrStdout(), "Maintenance created successfully.\n")
		fmt.Fprintf(cmd.OutOrStdout(), "Maintenance IDs: %v\n", maintenanceResponse.Result.MaintenanceIDs)
		fmt.Fprintf(cmd.OutOrStdout(), "Applied to %d host groups\n", len(groupIDs))
		fmt.Fprintf(cmd.OutOrStdout(), "Active from: %s\n", now.Format(time.RFC3339))
		fmt.Fprintf(cmd.OutOrStdout(), "Active until: %s\n", now.Add(time.Duration(maintenanceDuration)*time.Hour).Format(time.RFC3339))

		return nil
	},
}

func init() {
	// Add flags for maintenance-create-all command
	MaintenanceCreateAllCmd.Flags().StringVarP(&maintenanceName, "name", "n",
		fmt.Sprintf("Automatic Maintenance %s", time.Now().Format("2006-01-02")),
		"Name of the maintenance period")
	MaintenanceCreateAllCmd.Flags().StringVarP(&maintenanceDescription, "description", "d",
		"Automatically created maintenance period for all host groups",
		"Description of the maintenance period")
	MaintenanceCreateAllCmd.Flags().IntVarP(&maintenanceDuration, "duration", "t", 168,
		"Duration of the maintenance period in hours (default: 168 hours / 1 week)")
	MaintenanceCreateAllCmd.Flags().IntVarP(&maintenanceType, "type", "m",
		int(zabbix.MaintenanceWithDataCollection),
		"Type of maintenance: 0 (with data collection), 1 (no data collection)")

	// Add subcommands to maintenance command
	MaintenanceCmd.AddCommand(MaintenanceCreateAllCmd)
}

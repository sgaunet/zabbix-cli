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
	maintenanceType        zabbix.MaintenanceType
)

func init() {
	// Add flags for maintenance create command
	MaintenanceCreateCmd.Flags().StringVarP(&maintenanceName, "name", "n", "", "Name of the maintenance period (default: 'Maintenance <timestamp>')")
	MaintenanceCreateCmd.Flags().StringVarP(&maintenanceDescription, "description", "d", "", "Description of the maintenance period (default: 'Maintenance created by zabbix-cli on <timestamp>')")
	MaintenanceCreateCmd.Flags().IntVarP(&maintenanceDuration, "duration", "D", 1, "Duration of the maintenance period in hours")
	MaintenanceCreateCmd.Flags().IntVarP((*int)(&maintenanceType), "type", "t", int(zabbix.MaintenanceWithDataCollection), "Type of maintenance (0 - with data collection, 1 - without data collection)")
}

// MaintenanceCreateCmd represents the create subcommand that creates maintenance for all host groups
var MaintenanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a maintenance period for all host groups",
	Long:  `Create a one-time maintenance period for all host groups, active from now for the specified duration.`,
	Args:  cobra.ExactArgs(0),
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

		if maintenanceName == "" {
			maintenanceName = fmt.Sprintf("Maintenance %s", now.Format("2006-01-02 15:04:05"))
		}
		if maintenanceDescription == "" {
			maintenanceDescription = fmt.Sprintf("Maintenance created by zabbix-cli on %s", now.Format("2006-01-02 15:04:05"))
		}
		if maintenanceType == 0 {
			maintenanceType = zabbix.MaintenanceWithDataCollection
		}
		// Prepare maintenance creation request
		maintenanceRequest := zabbix.NewMaintenanceCreateRequest(
			zabbix.WithMaintenanceName(maintenanceName),
			zabbix.WithMaintenanceDescription(maintenanceDescription),
			zabbix.WithMaintenanceActiveTill(activeTill),
			zabbix.WithMaintenanceType(maintenanceType),
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

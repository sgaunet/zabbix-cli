package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var (
	maintenanceGetOutput    string
	maintenanceGetFormat    string
	maintenanceGetGroupIDs  string
	maintenanceGetHostIDs   string
	maintenanceGetIDs       string
	maintenanceGetLimit     int
	maintenanceGetSortField string
)

// MaintenanceGetCmd represents the maintenance get command
var MaintenanceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get maintenance periods",
	Long:  `Get maintenance periods from Zabbix with optional filtering.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if conf == nil {
			return fmt.Errorf("zabbix-cli is not configured properly")
		}

		ctx := context.Background()

		// Initialize the Zabbix client
		z := zabbix.New(conf.ZabbixUser, conf.ZabbixPassword, conf.ZabbixEndpoint)
		if err := z.Login(ctx); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		defer func() {
			if err := z.Logout(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "logout failed: %v\n", err)
			}
		}()

		// Prepare options for the maintenance.get request
		options := []zabbix.MaintenanceGetOption{
			zabbix.WithMaintenanceGetAuthToken(z.Auth()),
			zabbix.WithMaintenanceGetID(1), // Simple ID for CLI context
		}

		// Parse output parameter
		if maintenanceGetOutput != "" {
			options = append(options, zabbix.WithMaintenanceGetOutput(maintenanceGetOutput))
		}

		// Parse group IDs
		if maintenanceGetGroupIDs != "" {
			groupIDs := strings.Split(maintenanceGetGroupIDs, ",")
			options = append(options, zabbix.WithMaintenanceGetGroupIDs(groupIDs))
		}

		// Parse host IDs
		if maintenanceGetHostIDs != "" {
			hostIDs := strings.Split(maintenanceGetHostIDs, ",")
			options = append(options, zabbix.WithMaintenanceGetHostIDs(hostIDs))
		}

		// Parse maintenance IDs
		if maintenanceGetIDs != "" {
			maintenanceIDs := strings.Split(maintenanceGetIDs, ",")
			options = append(options, zabbix.WithMaintenanceGetMaintenanceIDs(maintenanceIDs))
		}

		// Set limit if provided
		if maintenanceGetLimit > 0 {
			options = append(options, zabbix.WithMaintenanceGetLimit(maintenanceGetLimit))
		}

		// Parse sort field
		if maintenanceGetSortField != "" {
			sortFields := strings.Split(maintenanceGetSortField, ",")
			options = append(options, zabbix.WithMaintenanceGetSortField(sortFields))
		}

		// Create and send the request
		request := zabbix.NewMaintenanceGetRequest(options...)
		response, err := z.MaintenanceGet(ctx, request)
		if err != nil {
			return fmt.Errorf("failed to get maintenance periods: %w", err)
		}

		// Output the results
		if maintenanceGetFormat == "json" {
			// Output in JSON format
			jsonBytes, err := json.MarshalIndent(response.Result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal result to JSON: %w", err)
			}
			fmt.Println(string(jsonBytes))
		} else {
			// Output in table format
			fmt.Println("Maintenance Periods:")
			fmt.Println("--------------------")

			for _, m := range response.Result {
				maintenanceType := "With Data Collection"
				if m.MaintenanceType == zabbix.MaintenanceNoDataCollection {
					maintenanceType = "No Data Collection"
				}

				// Format timestamps as human-readable dates
				activeSince := time.Unix(m.ActiveSince.Int64(), 0).Format("2006-01-02 15:04:05")
				activeTill := time.Unix(m.ActiveTill.Int64(), 0).Format("2006-01-02 15:04:05")

				fmt.Printf("ID: %s\n", m.MaintenanceID)
				fmt.Printf("Name: %s\n", m.Name)
				if m.Description != "" {
					fmt.Printf("Description: %s\n", m.Description)
				}
				fmt.Printf("Active Since: %s\n", activeSince)
				fmt.Printf("Active Till: %s\n", activeTill)
				fmt.Printf("Type: %s\n", maintenanceType)
				fmt.Printf("Host Groups: %d\n", len(m.GroupIDs))
				fmt.Printf("Hosts: %d\n", len(m.HostIDs))
				fmt.Println("--------------------")
			}

			fmt.Printf("Total maintenance periods: %d\n", len(response.Result))
		}

		return nil
	},
}

func init() {
	// Add flags for maintenance get command
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetOutput, "output", "o", "", "Fields to return (comma-separated)")
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetFormat, "format", "f", "table", "Output format: table or json")
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetGroupIDs, "groupids", "g", "", "Filter by host group IDs (comma-separated)")
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetHostIDs, "hostids", "H", "", "Filter by host IDs (comma-separated)")
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetIDs, "maintenanceids", "m", "", "Filter by maintenance IDs (comma-separated)")
	MaintenanceGetCmd.Flags().IntVarP(&maintenanceGetLimit, "limit", "l", 0, "Limit the number of results")
	MaintenanceGetCmd.Flags().StringVarP(&maintenanceGetSortField, "sort", "s", "", "Sort field(s) (comma-separated)")
}

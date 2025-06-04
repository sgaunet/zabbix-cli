package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

var hostGroupOutputFormat string

var GetHostGroupCmd = &cobra.Command{
	Use:   "hostgroup",
	Short: "Get Zabbix host groups",
	Long:  `Get Zabbix host groups. By default, displays results in a table. Use --output json for raw JSON.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// initConfig() should be called by cobra.OnInitialize on the RootCmd.
		// If conf is nil here, it means initConfig was not called or failed.
		if conf == nil {
			// This check can be made more robust depending on how initConfig signals failure.
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

		// Prepare request to get all host groups
		// Assuming zabbix.GenerateUniqueID() is available and similar to how Login handles it.
		// If not, we might need to make it available or use a simpler ID generation for CLI.
		// For now, let's assume ID=1 for simplicity in CLI context if GenerateUniqueID is not easily usable here.
		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth(z.Auth()),
			zabbix.WithHostGroupGetID(1), // Simplified ID for CLI context
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}

		if hostGroupOutputFormat == "json" {
			jsonOutput, err := json.MarshalIndent(response.Result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal host groups to JSON: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(jsonOutput))
		} else {
			// Table output
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tINTERNAL") // Add more columns as needed
			for _, hg := range response.Result {
				// Assuming HostGroup struct has GroupID, Name, and Internal fields.
				// Adjust according to the actual HostGroup struct definition.
				fmt.Fprintf(w, "%s\t%s\t%s\n", hg.GroupID, hg.Name, hg.Internal)
			}
			w.Flush()
		}

		return nil
	},
}

func init() {
	// Command registration will be handled by the parent command (get.go)
	GetHostGroupCmd.Flags().StringVarP(&hostGroupOutputFormat, "output", "o", "table", "Output format (table or json)")
}

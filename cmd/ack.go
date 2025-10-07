package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/spf13/cobra"
)

var ackSeverityFlag string
var ackDashboardName string

var ackCmd = &cobra.Command{
	Use:   "ack",
	Short: "acknowledge events",
	Long:  `acknowledge events`,
	Run: func(_ *cobra.Command, _ []string) {
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
		defer z.Logout(ctx) //nolint: errcheck

		var problemOptions []zabbix.GetProblemOption

		// Apply dashboard filters if specified
		if ackDashboardName != "" {
			filter := map[string]interface{}{
				"name": ackDashboardName,
			}
			dashboardReq := zabbix.NewDashboardGetRequest(
				zabbix.WithDashboardGetAuth(z.Auth()),
				zabbix.WithDashboardGetID(1),
				zabbix.WithDashboardGetOutput("extend"),
				zabbix.WithDashboardGetSelectPages("extend"),
				zabbix.WithDashboardGetFilter(filter),
			)

			dashboardResp, err := z.DashboardGet(ctx, dashboardReq)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get dashboard: %v\n", err.Error())
				os.Exit(1)
			}

			if len(dashboardResp.Result) == 0 {
				fmt.Fprintf(os.Stderr, "Dashboard '%s' not found\n", ackDashboardName)
				os.Exit(1)
			}

			dashboardFilters, err := zabbix.ParseProblemsWidgetFilters(&dashboardResp.Result[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse dashboard filters: %v\n", err.Error())
				os.Exit(1)
			}

			if len(dashboardFilters) == 0 {
				fmt.Fprintf(os.Stderr, "Warning: Dashboard '%s' has no 'problems' widget or no filters\n", ackDashboardName)
			}

			problemOptions = append(problemOptions, dashboardFilters...)
		}

		// Apply CLI severity flag (overrides dashboard severity if both are set)
		if ackSeverityFlag != "" {
			severityInt := zabbix.GetSeverityString(ackSeverityFlag)
			// ProblemParams.Severities expects []string of integer severities
			problemOptions = append(problemOptions, zabbix.GetProblemOptionSeverities([]string{fmt.Sprintf("%d", severityInt)}))
		}

		problems, err := z.GetProblems(ctx, problemOptions...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
			os.Exit(1)
		}

		for _, pb := range problems {
			fmt.Printf("Acknowledge problem %s\n", pb.Name)
			_, err = z.AcknowledgeEvents(ctx, []string{pb.EventID}, zabbix.WithActions(zabbix.AddMessage, zabbix.CloseProblem, zabbix.Acknowledge), zabbix.WithMessage("acknowledged from CLI"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to acknowledge problem %s: %v\n", pb.Name, err.Error())
				continue
			}
		}
	},
}

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

// loadDashboardFilters loads problem filters from a named dashboard.
// Returns the filters and an error if the dashboard cannot be loaded.
// Prints warnings to stderr for missing widgets or filters.
func loadDashboardFilters(ctx context.Context, z zabbix.Client, dashboardName string) ([]zabbix.GetProblemOption, error) {
	if dashboardName == "" {
		return nil, nil
	}

	filter := map[string]interface{}{
		"name": dashboardName,
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
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	if len(dashboardResp.Result) == 0 {
		return nil, fmt.Errorf("dashboard '%s' not found", dashboardName)
	}

	dashboardFilters, err := zabbix.ParseProblemsWidgetFilters(&dashboardResp.Result[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard filters: %w", err)
	}

	if len(dashboardFilters) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: Dashboard '%s' has no 'problems' widget or no filters\n", dashboardName)
	}

	return dashboardFilters, nil
}

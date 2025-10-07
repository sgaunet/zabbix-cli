package zabbix

import (
	"testing"
)

func TestExportRequestOptionDashboardsID(t *testing.T) {
	dashboardsID := []string{"1", "2", "3"}
	req := NewConfigurationExportRequest(ExportRequestOptionDashboardsID(dashboardsID))

	if len(req.Params.Options.DashboardsID) != 3 {
		t.Errorf("Expected 3 dashboard IDs, got %d", len(req.Params.Options.DashboardsID))
	}
	if req.Params.Options.DashboardsID[0] != "1" {
		t.Errorf("Expected first dashboard ID to be 1, got %s", req.Params.Options.DashboardsID[0])
	}
	if req.Params.Options.DashboardsID[1] != "2" {
		t.Errorf("Expected second dashboard ID to be 2, got %s", req.Params.Options.DashboardsID[1])
	}
}

func TestExportRequestOptionDashboardsIDWithYAMLFormat(t *testing.T) {
	dashboardsID := []string{"42"}
	req := NewConfigurationExportRequest(
		ExportRequestOptionDashboardsID(dashboardsID),
		ExportRequestOptionYAMLFormat(),
	)

	if req.Params.Options.DashboardsID[0] != "42" {
		t.Errorf("Expected dashboard ID to be 42, got %s", req.Params.Options.DashboardsID[0])
	}
	if req.Params.Format != "yaml" {
		t.Errorf("Expected format to be yaml, got %s", req.Params.Format)
	}
}

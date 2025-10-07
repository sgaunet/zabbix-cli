package zabbix

import (
	"testing"
)

func TestNewDashboardGetRequest(t *testing.T) {
	req := NewDashboardGetRequest()
	if req.JSONRPC != JSONRPC {
		t.Errorf("Expected JSONRPC to be %s, got %s", JSONRPC, req.JSONRPC)
	}
	if req.Method != "dashboard.get" {
		t.Errorf("Expected Method to be dashboard.get, got %s", req.Method)
	}
}

func TestWithDashboardGetDashboardIDs(t *testing.T) {
	ids := []string{"1", "2", "3"}
	req := NewDashboardGetRequest(WithDashboardGetDashboardIDs(ids))
	if len(req.Params.DashboardIDs) != 3 {
		t.Errorf("Expected 3 dashboard IDs, got %d", len(req.Params.DashboardIDs))
	}
	if req.Params.DashboardIDs[0] != "1" {
		t.Errorf("Expected first dashboard ID to be 1, got %s", req.Params.DashboardIDs[0])
	}
}

func TestWithDashboardGetSelectPages(t *testing.T) {
	req := NewDashboardGetRequest(WithDashboardGetSelectPages("extend"))
	if req.Params.SelectPages != "extend" {
		t.Errorf("Expected SelectPages to be extend, got %v", req.Params.SelectPages)
	}
}

func TestWithDashboardGetOutput(t *testing.T) {
	req := NewDashboardGetRequest(WithDashboardGetOutput("extend"))
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
}

func TestWithDashboardGetAuth(t *testing.T) {
	token := "test-auth-token"
	req := NewDashboardGetRequest(WithDashboardGetAuth(token))
	if req.Auth != token {
		t.Errorf("Expected Auth to be %s, got %s", token, req.Auth)
	}
}

func TestWithDashboardGetID(t *testing.T) {
	id := 123
	req := NewDashboardGetRequest(WithDashboardGetID(id))
	if req.ID != id {
		t.Errorf("Expected ID to be %d, got %d", id, req.ID)
	}
}

func TestWithDashboardGetFilter(t *testing.T) {
	filter := map[string]interface{}{
		"name": "test dashboard",
	}
	req := NewDashboardGetRequest(WithDashboardGetFilter(filter))
	if req.Params.Filter["name"] != "test dashboard" {
		t.Errorf("Expected filter name to be 'test dashboard', got %v", req.Params.Filter["name"])
	}
}

func TestNewGetAllDashboardsRequest(t *testing.T) {
	req := NewGetAllDashboardsRequest()
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
	if req.Params.SelectPages != "extend" {
		t.Errorf("Expected SelectPages to be extend, got %v", req.Params.SelectPages)
	}
}

func TestNewGetAllDashboardsRequestWithOverrides(t *testing.T) {
	token := "override-token"
	id := 999
	req := NewGetAllDashboardsRequest(
		WithDashboardGetAuth(token),
		WithDashboardGetID(id),
	)
	if req.Auth != token {
		t.Errorf("Expected Auth to be %s, got %s", token, req.Auth)
	}
	if req.ID != id {
		t.Errorf("Expected ID to be %d, got %d", id, req.ID)
	}
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
}

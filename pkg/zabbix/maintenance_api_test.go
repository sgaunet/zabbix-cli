package zabbix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMaintenanceCreate(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify request has expected fields
		if req.Method != "maintenance.create" {
			t.Errorf("Expected method 'maintenance.create', got '%s'", req.Method)
		}

		if req.Params.Name == "" {
			t.Error("Expected maintenance name to be set")
		}

		if len(req.Params.GroupIDs) == 0 {
			t.Error("Expected at least one host group ID")
		}

		if len(req.Params.TimePeriods) == 0 {
			t.Error("Expected at least one time period")
		}

		// Return success response
		resp := MaintenanceCreateResponse{
			JSONRPC: "2.0",
			Result: MaintenanceResponse{
				MaintenanceIDs: []string{"123"},
			},
			ID: req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create test client
	z := New("testuser", "testpassword", server.URL)
	z.auth = "testtoken"

	// Create test request
	timePeriod := []TimePeriod{
		{
			TimePeriodType: TimePeriodTypeOneTime,
			StartTime:      0,
			Period:         3600 * 24 * 7, // 1 week in seconds
		},
	}

	request := NewMaintenanceCreateRequest(
		WithMaintenanceName("Test Maintenance"),
		WithMaintenanceDescription("Test Description"),
		WithMaintenanceActiveSince(1717575000),
		WithMaintenanceActiveTill(1717575000 + 3600*24*7),
		WithMaintenanceType(MaintenanceWithDataCollection),
		WithMaintenanceTimePeriods(timePeriod),
		WithMaintenanceGroupIDs([]string{"1", "2"}),
		WithMaintenanceAuthToken("testtoken"),
		WithMaintenanceRequestID(1),
	)

	// Make the API call
	resp, err := z.MaintenanceCreate(context.Background(), request)
	if err != nil {
		t.Fatalf("MaintenanceCreate returned error: %v", err)
	}

	// Verify response
	if len(resp.Result.MaintenanceIDs) == 0 {
		t.Error("Expected maintenance IDs in response")
	}

	if resp.Result.MaintenanceIDs[0] != "123" {
		t.Errorf("Expected maintenance ID '123', got '%s'", resp.Result.MaintenanceIDs[0])
	}
}

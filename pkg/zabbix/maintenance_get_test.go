package zabbix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMaintenanceGet(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request body is correctly formed
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Check that the request method is correct
		if req.Method != "maintenance.get" {
			t.Errorf("Expected method 'maintenance.get', got '%s'", req.Method)
		}

		// Verify parameters
		if len(req.Params.GroupIDs) != 1 || req.Params.GroupIDs[0] != "4" {
			t.Errorf("Expected group ID 4, got %v", req.Params.GroupIDs)
		}

		// Return a mock response
		mockResponse := MaintenanceGetResponse{
			JSONRPC: "2.0",
			Result: []Maintenance{
				{
					MaintenanceID:   "1",
					Name:            "Test Maintenance",
					ActiveSince:     1609459200, // 2021-01-01 00:00:00
					ActiveTill:      1612137600, // 2021-02-01 00:00:00
					Description:     "Test maintenance description",
					MaintenanceType: MaintenanceWithDataCollection,
					GroupIDs:        []string{"4"},
					TimePeriods: []TimePeriod{
						{
							TimePeriodType: TimePeriodTypeOneTime,
							StartTime:      0,
							Period:         3600 * 24 * 7, // 1 week in seconds
						},
					},
				},
			},
			ID: req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	// Create a client that uses our test server
	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	// Create a maintenance.get request with options
	request := NewMaintenanceGetRequest(
		WithMaintenanceGetGroupIDs([]string{"4"}),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	// Call the MaintenanceGet method
	ctx := context.Background()
	response, err := client.MaintenanceGet(ctx, request)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the response fields
	if len(response.Result) != 1 {
		t.Fatalf("Expected 1 maintenance, got %d", len(response.Result))
	}

	maintenance := response.Result[0]
	if maintenance.MaintenanceID != "1" {
		t.Errorf("Expected maintenance ID '1', got '%s'", maintenance.MaintenanceID)
	}
	if maintenance.Name != "Test Maintenance" {
		t.Errorf("Expected name 'Test Maintenance', got '%s'", maintenance.Name)
	}
	if maintenance.Description != "Test maintenance description" {
		t.Errorf("Expected description 'Test maintenance description', got '%s'", maintenance.Description)
	}
	if len(maintenance.TimePeriods) != 1 {
		t.Fatalf("Expected 1 time period, got %d", len(maintenance.TimePeriods))
	}
}

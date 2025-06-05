package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/config"
	"github.com/spf13/cobra"
)

// MockZabbixAPI mocks the Zabbix API server for testing purposes
type MockZabbixAPI struct {
	server *httptest.Server
}

func NewMockZabbixAPI(t *testing.T) *MockZabbixAPI {
	mock := &MockZabbixAPI{}
	
	// Setup a test server
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		
		// Check request method and content type
		if r.Method != http.MethodPost {
			t.Fatalf("Expected POST request, got %s", r.Method)
		}
		
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			t.Fatalf("Expected Content-Type to contain application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		// Check which API method is being called and respond accordingly
		var req map[string]interface{}
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		
		method, ok := req["method"].(string)
		if !ok {
			t.Fatalf("Request doesn't contain method field")
		}
		
		w.Header().Set("Content-Type", "application/json")
		
		switch method {
		case "user.login":
			// Return a successful login response
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  "mock-session-token",
				"id":      req["id"],
			}
			json.NewEncoder(w).Encode(response)
			
		case "hostgroup.get":
			// Return some mock host groups
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"result": []map[string]interface{}{
					{
						"groupid":  "1",
						"name":     "Linux servers",
						"internal": "0",
					},
					{
						"groupid":  "2",
						"name":     "Windows servers",
						"internal": "0",
					},
				},
				"id": req["id"],
			}
			json.NewEncoder(w).Encode(response)
			
		case "maintenance.create":
			// Check that the maintenance create request contains expected fields
			params, ok := req["params"].(map[string]interface{})
			if !ok {
				t.Fatalf("Request doesn't contain params field")
			}
			
			_, nameOK := params["name"]
			groupIDs, groupsOK := params["groupids"].([]interface{})
			timePeriods, periodsOK := params["timeperiods"].([]interface{})
			
			if !nameOK || !groupsOK || !periodsOK {
				t.Fatalf("Maintenance create request missing required fields")
			}
			
			if len(groupIDs) == 0 {
				t.Fatalf("Expected at least one host group ID")
			}
			
			if len(timePeriods) == 0 {
				t.Fatalf("Expected at least one time period")
			}
			
			// Return a successful maintenance create response
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"maintenanceids": []string{"123"},
				},
				"id": req["id"],
			}
			json.NewEncoder(w).Encode(response)
			
		case "user.logout":
			// Return a successful logout response
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"result":  true,
				"id":      req["id"],
			}
			json.NewEncoder(w).Encode(response)
			
		default:
			t.Fatalf("Unexpected API method: %s", method)
		}
	}))
	
	return mock
}

func (m *MockZabbixAPI) Close() {
	m.server.Close()
}

func (m *MockZabbixAPI) URL() string {
	return m.server.URL
}

func TestMaintenanceCreateAllCmd(t *testing.T) {
	// Skip integration tests when running in CI
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Setup mock API
	mockAPI := NewMockZabbixAPI(t)
	defer mockAPI.Close()
	
	// Setup test configuration - IMPORTANT: Use the mock URL
	oldConf := conf // Save any existing config
	conf = &config.Config{
		ZabbixEndpoint: mockAPI.URL(), // This must match the mock server URL
		ZabbixUser:     "admin",
		ZabbixPassword: "zabbix",
	}
	defer func() { conf = oldConf }() // Restore old config after test
	
	// Setup command for testing
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(MaintenanceCreateAllCmd)
	output := &bytes.Buffer{}
	cmd.SetOut(output)
	
	// Execute the command
	cmd.SetArgs([]string{"create-all", "--name", "Test Maintenance"})
	err := cmd.Execute()
	
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}
	
	// Check the command output
	if !strings.Contains(output.String(), "Maintenance created successfully") {
		t.Errorf("Expected output to contain 'Maintenance created successfully', got '%s'", output.String())
	}
	
	if !strings.Contains(output.String(), "Maintenance IDs: [123]") {
		t.Errorf("Expected output to contain 'Maintenance IDs: [123]', got '%s'", output.String())
	}
	
	if !strings.Contains(output.String(), "Applied to 2 host groups") {
		t.Errorf("Expected output to contain 'Applied to 2 host groups', got '%s'", output.String())
	}
}

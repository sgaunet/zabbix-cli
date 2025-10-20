package zabbix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
			JSONRPC: JSONRPC,
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

func TestMaintenanceCreateWithTagsEvalType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		tagsEvalType  int
		expectedValue int
	}{
		{
			name:          "TagsEvalType AND (0)",
			tagsEvalType:  int(TagsEvalTypeAnd),
			expectedValue: 0,
		},
		{
			name:          "TagsEvalType OR (1)",
			tagsEvalType:  int(TagsEvalTypeOr),
			expectedValue: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req MaintenanceCreateRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("Failed to decode request: %v", err)
				}

				// Verify tags_evaltype is set correctly
				if req.Params.TagsEvalType != tt.expectedValue {
					t.Errorf("Expected tags_evaltype %d, got %d", tt.expectedValue, req.Params.TagsEvalType)
				}

				// Return success response
				resp := MaintenanceCreateResponse{
					JSONRPC: JSONRPC,
					Result: MaintenanceResponse{
						MaintenanceIDs: []string{"456"},
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

			// Create test request with tags and tags_evaltype
			timePeriod := []TimePeriod{
				{
					TimePeriodType: TimePeriodTypeOneTime,
					StartTime:      0,
					Period:         3600,
				},
			}

			tags := []ProblemTag{
				{Tag: "environment", Value: "production"},
				{Tag: "application", Value: "web"},
			}

			request := NewMaintenanceCreateRequest(
				WithMaintenanceName("Test Maintenance with Tags"),
				WithMaintenanceActiveSince(1717575000),
				WithMaintenanceActiveTill(1717575000+3600),
				WithMaintenanceType(MaintenanceWithDataCollection),
				WithMaintenanceTimePeriods(timePeriod),
				WithMaintenanceHostIDs([]string{"1"}),
				WithMaintenanceTags(tags),
				WithMaintenanceTagsEvalType(tt.tagsEvalType),
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
		})
	}
}

func TestTagsEvalTypeConstants(t *testing.T) {
	t.Parallel()

	// Verify TagsEvalType constants match Zabbix API 7.2 spec
	if TagsEvalTypeAnd != 0 {
		t.Errorf("TagsEvalTypeAnd should be 0, got %d", TagsEvalTypeAnd)
	}

	if TagsEvalTypeOr != 1 {
		t.Errorf("TagsEvalTypeOr should be 1, got %d", TagsEvalTypeOr)
	}
}

func TestMaintenanceCreateRequestJSONMarshaling(t *testing.T) {
	t.Parallel()

	// Create a maintenance request with tags_evaltype
	timePeriod := []TimePeriod{
		{
			TimePeriodType: TimePeriodTypeOneTime,
			StartTime:      0,
			Period:         3600,
		},
	}

	tags := []ProblemTag{
		{Tag: "env", Value: "prod"},
	}

	request := NewMaintenanceCreateRequest(
		WithMaintenanceName("Test"),
		WithMaintenanceActiveSince(1717575000),
		WithMaintenanceActiveTill(1717575000+3600),
		WithMaintenanceType(MaintenanceWithDataCollection),
		WithMaintenanceTimePeriods(timePeriod),
		WithMaintenanceHostIDs([]string{"1"}),
		WithMaintenanceTags(tags),
		WithMaintenanceTagsEvalType(1), // OR
		WithMaintenanceAuthToken("token"),
		WithMaintenanceRequestID(1),
	)

	// Marshal to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Verify tags_evaltype is present in JSON
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "tags_evaltype") {
		t.Error("Expected tags_evaltype to be present in JSON")
	}

	// Unmarshal back and verify
	var unmarshaled MaintenanceCreateRequest
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if unmarshaled.Params.TagsEvalType != 1 {
		t.Errorf("Expected tags_evaltype 1 after unmarshal, got %d", unmarshaled.Params.TagsEvalType)
	}
}

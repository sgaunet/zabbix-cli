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
		tagsEvalType  TagsEvalType
		expectedValue TagsEvalType
	}{
		{
			name:          "TagsEvalType AND (0)",
			tagsEvalType:  TagsEvalTypeAnd,
			expectedValue: TagsEvalTypeAnd,
		},
		{
			name:          "TagsEvalType OR (1)",
			tagsEvalType:  TagsEvalTypeOr,
			expectedValue: TagsEvalTypeOr,
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
		WithMaintenanceTagsEvalType(TagsEvalTypeOr),
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

// TestMaintenanceTagsEvalTypeUnmarshalCompat tests unmarshaling tags_evaltype from both string and int
// This ensures compatibility with both Zabbix 6 (returns strings) and Zabbix 7 (may return ints)
func TestMaintenanceTagsEvalTypeUnmarshalCompat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		jsonData string
		expected TagsEvalType
	}{
		{
			name:     "tags_evaltype as string '0' (Zabbix 6)",
			jsonData: `{"tags_evaltype":"0"}`,
			expected: TagsEvalTypeAnd,
		},
		{
			name:     "tags_evaltype as string '1' (Zabbix 6)",
			jsonData: `{"tags_evaltype":"1"}`,
			expected: TagsEvalTypeOr,
		},
		{
			name:     "tags_evaltype as int 0 (Zabbix 7)",
			jsonData: `{"tags_evaltype":0}`,
			expected: TagsEvalTypeAnd,
		},
		{
			name:     "tags_evaltype as int 1 (Zabbix 7)",
			jsonData: `{"tags_evaltype":1}`,
			expected: TagsEvalTypeOr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var maint struct {
				TagsEvalType TagsEvalType `json:"tags_evaltype"`
			}

			err := json.Unmarshal([]byte(tt.jsonData), &maint)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if maint.TagsEvalType != tt.expected {
				t.Errorf("Expected TagsEvalType %d, got %d", tt.expected, maint.TagsEvalType)
			}
		})
	}
}

// TestMaintenanceGetResponseZabbix6Compat tests unmarshaling a real Zabbix 6 response
func TestMaintenanceGetResponseZabbix6Compat(t *testing.T) {
	t.Parallel()

	// This is the actual JSON response from Zabbix 6 that was failing
	zabbix6Response := `{
		"jsonrpc":"2.0",
		"result":[
			{
				"maintenanceid":"50",
				"name":"Test maintenance",
				"maintenance_type":"0",
				"description":"Maintenance created by zabbix-cli on 2025-11-13 13:55:55",
				"active_since":"0",
				"active_till":"1763211300",
				"tags_evaltype":"0"
			},
			{
				"maintenanceid":"47",
				"name":"Maintenance",
				"maintenance_type":"0",
				"description":"",
				"active_since":"1722621600",
				"active_till":"1724349600",
				"tags_evaltype":"0"
			}
		],
		"id":1
	}`

	var response struct {
		JSONRPC string        `json:"jsonrpc"`
		Result  []Maintenance `json:"result"`
		ID      int           `json:"id"`
	}

	err := json.Unmarshal([]byte(zabbix6Response), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal Zabbix 6 response: %v", err)
	}

	if len(response.Result) != 2 {
		t.Fatalf("Expected 2 maintenance periods, got %d", len(response.Result))
	}

	// Verify first maintenance
	maint1 := response.Result[0]
	if maint1.MaintenanceID != "50" {
		t.Errorf("Expected maintenanceid '50', got '%s'", maint1.MaintenanceID)
	}
	if maint1.Name != "Test maintenance" {
		t.Errorf("Expected name 'Test maintenance', got '%s'", maint1.Name)
	}
	if maint1.TagsEvalType != TagsEvalTypeAnd {
		t.Errorf("Expected tags_evaltype AND (0), got %d", maint1.TagsEvalType)
	}
	if maint1.MaintenanceType != MaintenanceWithDataCollection {
		t.Errorf("Expected maintenance_type 0, got %d", maint1.MaintenanceType)
	}

	// Verify second maintenance
	maint2 := response.Result[1]
	if maint2.MaintenanceID != "47" {
		t.Errorf("Expected maintenanceid '47', got '%s'", maint2.MaintenanceID)
	}
	if maint2.TagsEvalType != TagsEvalTypeAnd {
		t.Errorf("Expected tags_evaltype AND (0), got %d", maint2.TagsEvalType)
	}
}

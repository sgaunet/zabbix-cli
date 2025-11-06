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
			JSONRPC: JSONRPC,
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

func TestMaintenanceGetWithSelectExtend(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify Select parameters are "extend"
		if req.Params.SelectGroups != "extend" {
			t.Errorf("Expected SelectGroups 'extend', got %v", req.Params.SelectGroups)
		}
		if req.Params.SelectHosts != "extend" {
			t.Errorf("Expected SelectHosts 'extend', got %v", req.Params.SelectHosts)
		}
		if req.Params.SelectTimePeriods != "extend" {
			t.Errorf("Expected SelectTimePeriods 'extend', got %v", req.Params.SelectTimePeriods)
		}
		if req.Params.SelectTags != "extend" {
			t.Errorf("Expected SelectTags 'extend', got %v", req.Params.SelectTags)
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetSelectGroups("extend"),
		WithMaintenanceGetSelectHosts("extend"),
		WithMaintenanceGetSelectTimePeriods("extend"),
		WithMaintenanceGetSelectTags("extend"),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMaintenanceGetWithSelectArray(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify Select parameters are arrays
		groupsArray, ok := req.Params.SelectGroups.([]any)
		if !ok {
			t.Errorf("Expected SelectGroups to be array, got %T", req.Params.SelectGroups)
		} else if len(groupsArray) != 2 {
			t.Errorf("Expected SelectGroups array length 2, got %d", len(groupsArray))
		}

		tagsArray, ok := req.Params.SelectTags.([]any)
		if !ok {
			t.Errorf("Expected SelectTags to be array, got %T", req.Params.SelectTags)
		} else if len(tagsArray) != 2 {
			t.Errorf("Expected SelectTags array length 2, got %d", len(tagsArray))
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetSelectGroups([]string{"groupid", "name"}),
		WithMaintenanceGetSelectTags([]string{"tag", "value"}),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMaintenanceGetWithLimitSelects(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify LimitSelects is set
		if req.Params.LimitSelects != 10 {
			t.Errorf("Expected LimitSelects 10, got %d", req.Params.LimitSelects)
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetLimitSelects(10),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMaintenanceGetWithFilter(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify Filter is set
		if req.Params.Filter == nil {
			t.Error("Expected Filter to be set")
		} else {
			name, ok := req.Params.Filter["name"]
			if !ok {
				t.Error("Expected Filter to have 'name' key")
			} else if name != "Test Maintenance" {
				t.Errorf("Expected Filter name 'Test Maintenance', got %v", name)
			}
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	filter := map[string]any{
		"name": "Test Maintenance",
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetFilter(filter),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMaintenanceGetWithSearch(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify Search parameters are set
		if req.Params.Search == nil {
			t.Error("Expected Search to be set")
		}
		if !req.Params.SearchByAny {
			t.Error("Expected SearchByAny to be true")
		}
		if !req.Params.SearchWildcardsEnabled {
			t.Error("Expected SearchWildcardsEnabled to be true")
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	search := map[string]any{
		"name": "maintenance*",
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetSearch(search),
		WithMaintenanceGetSearchByAny(true),
		WithMaintenanceGetSearchWildcardsEnabled(true),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMaintenanceGetWithCommonGetParams(t *testing.T) {
	t.Parallel()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MaintenanceGetRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Verify CommonGetParams flags are set
		if !req.Params.StartSearch {
			t.Error("Expected StartSearch to be true")
		}
		if !req.Params.ExcludeSearch {
			t.Error("Expected ExcludeSearch to be true")
		}
		if !req.Params.PreserveKeys {
			t.Error("Expected PreserveKeys to be true")
		}
		if !req.Params.CountOutput {
			t.Error("Expected CountOutput to be true")
		}
		if !req.Params.Editable {
			t.Error("Expected Editable to be true")
		}

		mockResponse := MaintenanceGetResponse{
			JSONRPC: JSONRPC,
			Result:  []Maintenance{},
			ID:      req.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	client := &Client{
		APIEndpoint: ts.URL,
		client:      &http.Client{},
	}

	request := NewMaintenanceGetRequest(
		WithMaintenanceGetStartSearch(true),
		WithMaintenanceGetExcludeSearch(true),
		WithMaintenanceGetPreserveKeys(true),
		WithMaintenanceGetCountOutput(true),
		WithMaintenanceGetEditable(true),
		WithMaintenanceGetAuthToken("test-token"),
		WithMaintenanceGetID(1),
	)

	ctx := context.Background()
	_, err := client.MaintenanceGet(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaintenanceDelete(t *testing.T) {
	// Setup a test server to mock the Zabbix API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for correct HTTP method
		assert.Equal(t, http.MethodPost, r.Method)

		// Unmarshal the request body
		var request MaintenanceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify the request is properly formed
		assert.Equal(t, JSONRPC, request.JSONRPC)
		assert.Equal(t, "maintenance.delete", request.Method)
		assert.Equal(t, "dummy-auth-token", request.Auth)
		assert.NotEmpty(t, request.ID)
		assert.Equal(t, []string{"1", "2"}, request.Params)

		// Create a mock response
		response := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Result:  MaintenanceDeleteResult{
				MaintenanceIDs: []string{"1", "2"},
			},
			ID:      request.ID,
		}

		// Send the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a client pointing to the test server
	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Create a delete request with test data
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteIDs([]string{"1", "2"}),
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	// Call the method under test
	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	// Verify the results
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []string{"1", "2"}, response.Result.MaintenanceIDs)
	assert.Equal(t, JSONRPC, response.JSONRPC)
	assert.Equal(t, request.ID, response.ID)
	assert.Nil(t, response.Error)
}

func TestMaintenanceDeleteWithError(t *testing.T) {
	// Setup a test server to return an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create an error response
		errorResponse := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Error: &Error{
				Code:    42,
				Message: "Test error",
				Data:    "Maintenance periods not found",
			},
			ID: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	// Create a client pointing to the test server
	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Create a delete request with test data
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteID("99"),
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	// Call the method under test
	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	// Verify the results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "Zabbix API error (Code: 42): Test error - Maintenance periods not found", err.Error())
}

func TestMaintenanceDeleteSingleID(t *testing.T) {
	t.Parallel()

	// Setup a test server to verify single ID deletion
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request MaintenanceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify single ID in array
		assert.Equal(t, []string{"123"}, request.Params)
		assert.Len(t, request.Params, 1)

		response := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Result: MaintenanceDeleteResult{
				MaintenanceIDs: []string{"123"},
			},
			ID: request.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Use WithMaintenanceDeleteID for single ID
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteID("123"),
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []string{"123"}, response.Result.MaintenanceIDs)
}

func TestMaintenanceDeleteMultipleIDsIncremental(t *testing.T) {
	t.Parallel()

	// Setup a test server to verify multiple IDs added incrementally
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request MaintenanceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify all IDs were added
		assert.Equal(t, []string{"1", "2", "3"}, request.Params)
		assert.Len(t, request.Params, 3)

		response := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Result: MaintenanceDeleteResult{
				MaintenanceIDs: request.Params,
			},
			ID: request.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Add IDs incrementally using WithMaintenanceDeleteID multiple times
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteID("1"),
		WithMaintenanceDeleteID("2"),
		WithMaintenanceDeleteID("3"),
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []string{"1", "2", "3"}, response.Result.MaintenanceIDs)
}

func TestMaintenanceDeleteEmptyParams(t *testing.T) {
	t.Parallel()

	// Setup a test server to handle empty params
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request MaintenanceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify params is empty array, not null
		assert.NotNil(t, request.Params)
		assert.Len(t, request.Params, 0)

		// Zabbix would return an error for empty params, simulate this
		errorResponse := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Array is empty",
			},
			ID: request.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Create request with no IDs
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestMaintenanceDeleteJSONMarshaling(t *testing.T) {
	t.Parallel()

	// Test that the request marshals correctly to JSON
	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteIDs([]string{"100", "200", "300"}),
		WithMaintenanceDeleteAuthToken("test-token"),
	)
	request.ID = 1 // Set fixed ID for predictable JSON

	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	// Unmarshal to verify structure
	var unmarshaled map[string]any
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, JSONRPC, unmarshaled["jsonrpc"])
	assert.Equal(t, "maintenance.delete", unmarshaled["method"])
	assert.Equal(t, "test-token", unmarshaled["auth"])
	assert.Equal(t, float64(1), unmarshaled["id"])

	// Verify params is an array
	params, ok := unmarshaled["params"].([]any)
	assert.True(t, ok, "params should be an array")
	assert.Len(t, params, 3)
	assert.Equal(t, "100", params[0])
	assert.Equal(t, "200", params[1])
	assert.Equal(t, "300", params[2])
}

func TestMaintenanceDeleteLargeNumberOfIDs(t *testing.T) {
	t.Parallel()

	// Setup a test server to handle many IDs
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request MaintenanceDeleteRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify we received all 50 IDs
		assert.Len(t, request.Params, 50)

		response := MaintenanceDeleteResponse{
			JSONRPC: JSONRPC,
			Result: MaintenanceDeleteResult{
				MaintenanceIDs: request.Params,
			},
			ID: request.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "dummy-auth-token",
		client:      &http.Client{},
	}

	// Generate 50 IDs
	ids := make([]string, 50)
	for i := 0; i < 50; i++ {
		ids[i] = fmt.Sprintf("%d", i+1)
	}

	request := NewMaintenanceDeleteRequest(
		WithMaintenanceDeleteIDs(ids),
		WithMaintenanceDeleteAuthToken("dummy-auth-token"),
	)

	ctx := context.Background()
	response, err := client.MaintenanceDelete(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Result.MaintenanceIDs, 50)
}

func TestMaintenanceDeleteResponseUnmarshaling(t *testing.T) {
	t.Parallel()

	// Test that the response unmarshals correctly from JSON
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"maintenanceids": ["10", "20", "30"]
		},
		"id": 1
	}`

	var response MaintenanceDeleteResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	assert.NoError(t, err)

	assert.Equal(t, JSONRPC, response.JSONRPC)
	assert.Equal(t, 1, response.ID)
	assert.Nil(t, response.Error)
	assert.Equal(t, []string{"10", "20", "30"}, response.Result.MaintenanceIDs)
	assert.Len(t, response.Result.MaintenanceIDs, 3)
}

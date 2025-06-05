package zabbix

import (
	"context"
	"encoding/json"
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

package zabbix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHostGroupCreateSingle tests creating a single host group
func TestHostGroupCreateSingle(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var request HostGroupCreateRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify request is properly formed
		assert.Equal(t, JSONRPC, request.JSONRPC)
		assert.Equal(t, "hostgroup.create", request.Method)
		assert.Equal(t, "test-auth-token", request.Auth)
		assert.NotEmpty(t, request.ID)

		// Verify single group in params array
		assert.Equal(t, 1, len(request.Params))
		assert.Equal(t, "Test Group", request.Params[0].Name)

		// Create a mock response
		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Result: HostGroupCreateResponseData{
				GroupIDs: []string{"100"},
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
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	request := NewHostGroupCreateRequest(
		[]string{"Test Group"},
		WithHostGroupCreateAuth("test-auth-token"),
		WithHostGroupCreateID(1),
	)

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Result.GroupIDs))
	assert.Equal(t, "100", response.Result.GroupIDs[0])
	assert.Nil(t, response.Error)
}

// TestHostGroupCreateMultiple tests creating multiple host groups
func TestHostGroupCreateMultiple(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupCreateRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify multiple groups in params array
		assert.Equal(t, 3, len(request.Params))
		assert.Equal(t, "Group 1", request.Params[0].Name)
		assert.Equal(t, "Group 2", request.Params[1].Name)
		assert.Equal(t, "Group 3", request.Params[2].Name)

		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Result: HostGroupCreateResponseData{
				GroupIDs: []string{"101", "102", "103"},
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
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	request := NewHostGroupCreateRequest(
		[]string{"Group 1", "Group 2", "Group 3"},
		WithHostGroupCreateAuth("test-auth-token"),
		WithHostGroupCreateID(1),
	)

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Result.GroupIDs))
	assert.Equal(t, "101", response.Result.GroupIDs[0])
	assert.Equal(t, "102", response.Result.GroupIDs[1])
	assert.Equal(t, "103", response.Result.GroupIDs[2])
}

// TestHostGroupCreateDuplicateName tests error handling for duplicate group names
func TestHostGroupCreateDuplicateName(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate Zabbix API error for duplicate name
		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Error: &Error{
				Code:    -32500,
				Message: "Application error.",
				Data:    "Host group \"Test Group\" already exists.",
			},
			ID: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	request := NewHostGroupCreateRequest(
		[]string{"Test Group"},
		WithHostGroupCreateAuth("test-auth-token"),
		WithHostGroupCreateID(1),
	)

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	zbxErr, ok := err.(*Error)
	assert.True(t, ok, "error should be of type *zabbix.Error")
	assert.Equal(t, -32500, zbxErr.Code)
	assert.Contains(t, zbxErr.Data, "already exists")
}

// TestHostGroupCreateEmptyNames tests that nil is returned for empty group names
func TestHostGroupCreateEmptyNames(t *testing.T) {
	t.Parallel()

	// NewHostGroupCreateRequest should return nil for empty group names
	request := NewHostGroupCreateRequest([]string{})
	assert.Nil(t, request)
}

// TestHostGroupCreateJSONMarshaling tests that request marshals correctly to JSON
func TestHostGroupCreateJSONMarshaling(t *testing.T) {
	t.Parallel()

	request := NewHostGroupCreateRequest(
		[]string{"Group A", "Group B"},
		WithHostGroupCreateAuth("test-token"),
		WithHostGroupCreateID(5),
	)

	jsonData, err := json.Marshal(request)
	assert.NoError(t, err)

	// Unmarshal to verify structure
	var unmarshaled map[string]any
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Verify JSON structure
	assert.Equal(t, JSONRPC, unmarshaled["jsonrpc"])
	assert.Equal(t, "hostgroup.create", unmarshaled["method"])
	assert.Equal(t, "test-token", unmarshaled["auth"])
	assert.Equal(t, float64(5), unmarshaled["id"])

	// Verify params is an array
	params, ok := unmarshaled["params"].([]any)
	assert.True(t, ok, "params should be an array")
	assert.Equal(t, 2, len(params))

	// Verify first group
	group1, ok := params[0].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "Group A", group1["name"])

	// Verify second group
	group2, ok := params[1].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "Group B", group2["name"])
}

// TestHostGroupCreateResponseUnmarshaling tests correct response unmarshaling
func TestHostGroupCreateResponseUnmarshaling(t *testing.T) {
	t.Parallel()

	// Test that the response unmarshals correctly from JSON
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": {
			"groupids": ["10", "20", "30"]
		},
		"id": 1
	}`

	var response HostGroupCreateResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	assert.NoError(t, err)

	assert.Equal(t, JSONRPC, response.JSONRPC)
	assert.Equal(t, 1, response.ID)
	assert.Nil(t, response.Error)
	assert.Equal(t, 3, len(response.Result.GroupIDs))
	assert.Equal(t, "10", response.Result.GroupIDs[0])
	assert.Equal(t, "20", response.Result.GroupIDs[1])
	assert.Equal(t, "30", response.Result.GroupIDs[2])
}

// TestHostGroupCreateWithInvalidParams tests error handling for invalid parameters
func TestHostGroupCreateWithInvalidParams(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate Zabbix API error for invalid params
		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Error: &Error{
				Code:    -32602,
				Message: "Invalid params.",
				Data:    "Invalid parameter \"/1/name\": cannot be empty.",
			},
			ID: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	// Create a request with an empty name (simulating invalid input)
	request := &HostGroupCreateRequest{
		JSONRPC: JSONRPC,
		Method:  "hostgroup.create",
		Params:  []HostGroup{{Name: ""}}, // Empty name
		Auth:    "test-auth-token",
		ID:      1,
	}

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	zbxErr, ok := err.(*Error)
	assert.True(t, ok, "error should be of type *zabbix.Error")
	assert.Equal(t, -32602, zbxErr.Code)
	assert.Contains(t, zbxErr.Data, "cannot be empty")
}

// TestHostGroupCreateGroupIDsAsStrings tests that all group IDs are strings
func TestHostGroupCreateGroupIDsAsStrings(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Result: HostGroupCreateResponseData{
				GroupIDs: []string{"1", "2", "3", "100"},
			},
			ID: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := Client{
		APIEndpoint: server.URL,
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	request := NewHostGroupCreateRequest(
		[]string{"A", "B", "C", "D"},
		WithHostGroupCreateAuth("test-auth-token"),
		WithHostGroupCreateID(1),
	)

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Verify all IDs are strings
	for i, id := range response.Result.GroupIDs {
		assert.IsType(t, "", id, "GroupID at index %d should be a string", i)
	}
}

// TestHostGroupCreateLargeNumberOfGroups tests creating many groups at once
func TestHostGroupCreateLargeNumberOfGroups(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupCreateRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify we received 20 groups
		assert.Equal(t, 20, len(request.Params))

		// Generate 20 group IDs
		groupIDs := make([]string, 20)
		for i := 0; i < 20; i++ {
			groupIDs[i] = string(rune(200 + i))
		}

		response := HostGroupCreateResponse{
			JSONRPC: JSONRPC,
			Result: HostGroupCreateResponseData{
				GroupIDs: groupIDs,
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
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	// Generate 20 group names
	groupNames := make([]string, 20)
	for i := 0; i < 20; i++ {
		groupNames[i] = string(rune('A' + i)) + " Group"
	}

	request := NewHostGroupCreateRequest(
		groupNames,
		WithHostGroupCreateAuth("test-auth-token"),
		WithHostGroupCreateID(1),
	)

	ctx := context.Background()
	response, err := client.HostGroupCreate(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 20, len(response.Result.GroupIDs))
}

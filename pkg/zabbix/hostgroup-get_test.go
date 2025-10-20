package zabbix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHostGroupGetBasic tests basic hostgroup.get functionality
func TestHostGroupGetBasic(t *testing.T) {
	t.Parallel()

	// Setup a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify request is properly formed
		assert.Equal(t, JSONRPC, request.JSONRPC)
		assert.Equal(t, "hostgroup.get", request.Method)
		assert.Equal(t, "test-auth-token", request.Auth)
		assert.NotEmpty(t, request.ID)

		// Create a mock response with complete HostGroup structure
		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{
					GroupID:  "1",
					Name:     "Linux servers",
					Flags:    "0",
					Internal: "0",
					UUID:     "dc84c4f5e0e14d79a8a97f8b8f5e5b1f",
				},
				{
					GroupID:  "2",
					Name:     "Windows servers",
					Flags:    "0",
					Internal: "0",
					UUID:     "e95f3e2a1d2f4c8a9b7c6d5e4f3a2b1c",
				},
			},
			ID: request.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a client pointing to the test server
	client := Client{
		APIEndpoint: server.URL,
		auth:        "test-auth-token",
		client:      &http.Client{},
	}

	// Create a request
	request := NewHostGroupGetRequest(
		WithHostGroupGetAuth("test-auth-token"),
		WithHostGroupGetOutput("extend"),
		WithHostGroupGetID(1),
	)

	// Call the method under test
	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)

	// Verify the results
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Result))

	// Verify all fields are strings
	assert.Equal(t, "1", response.Result[0].GroupID)
	assert.Equal(t, "Linux servers", response.Result[0].Name)
	assert.Equal(t, "0", response.Result[0].Flags)
	assert.Equal(t, "0", response.Result[0].Internal)
	assert.Equal(t, "dc84c4f5e0e14d79a8a97f8b8f5e5b1f", response.Result[0].UUID)
	assert.Nil(t, response.Error)
}

// TestHostGroupGetWithFilter tests filtering by group name
func TestHostGroupGetWithFilter(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify filter is set
		assert.NotNil(t, request.Params.Filter)
		name, ok := request.Params.Filter["name"]
		assert.True(t, ok)
		assert.Equal(t, "Linux servers", name)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{
					GroupID: "1",
					Name:    "Linux servers",
				},
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetFilter(map[string]any{"name": "Linux servers"}),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, len(response.Result))
	assert.Equal(t, "Linux servers", response.Result[0].Name)
}

// TestHostGroupGetWithSelectExtend tests select options with "extend"
func TestHostGroupGetWithSelectExtend(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify Select parameters are "extend"
		assert.Equal(t, "extend", request.Params.SelectHosts)
		assert.Equal(t, "extend", request.Params.SelectGroupDiscoveries)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result:  []HostGroup{},
			ID:      request.ID,
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetSelectHosts("extend"),
		WithHostGroupGetSelectGroupDiscoveries("extend"),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	_, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
}

// TestHostGroupGetWithSelectArray tests select options with field arrays
func TestHostGroupGetWithSelectArray(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify Select parameters are arrays
		hostsArray, ok := request.Params.SelectHosts.([]any)
		assert.True(t, ok, "SelectHosts should be array")
		assert.Equal(t, 2, len(hostsArray))

		itemsArray, ok := request.Params.SelectItems.([]any)
		assert.True(t, ok, "SelectItems should be array")
		assert.Equal(t, 2, len(itemsArray))

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result:  []HostGroup{},
			ID:      request.ID,
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetSelectHosts([]string{"hostid", "name"}),
		WithHostGroupGetSelectItems([]string{"itemid", "key_"}),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	_, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
}

// TestHostGroupGetWithLimitSelects tests limitSelects parameter
func TestHostGroupGetWithLimitSelects(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify LimitSelects is set
		assert.Equal(t, 10, request.Params.LimitSelects)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result:  []HostGroup{},
			ID:      request.ID,
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetLimitSelects(10),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	_, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
}

// TestHostGroupGetWithPagination tests limit and offset for pagination
func TestHostGroupGetWithPagination(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify pagination parameters
		assert.Equal(t, 10, request.Params.Limit)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{GroupID: "1", Name: "Group 1"},
				{GroupID: "2", Name: "Group 2"},
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetLimit(10),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.LessOrEqual(t, len(response.Result), 10)
}

// TestHostGroupGetWithSorting tests sortfield and sortorder parameters
func TestHostGroupGetWithSorting(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify sorting parameters
		assert.Equal(t, []string{"name"}, request.Params.SortField)
		assert.Equal(t, []string{"ASC"}, request.Params.SortOrder)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{GroupID: "1", Name: "A Group"},
				{GroupID: "2", Name: "B Group"},
				{GroupID: "3", Name: "C Group"},
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetSortField([]string{"name"}),
		WithHostGroupGetSortOrder([]string{"ASC"}),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Result))
	// Verify results are sorted
	assert.Equal(t, "A Group", response.Result[0].Name)
	assert.Equal(t, "B Group", response.Result[1].Name)
	assert.Equal(t, "C Group", response.Result[2].Name)
}

// TestHostGroupGetWithSearch tests search parameters
func TestHostGroupGetWithSearch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify Search parameters
		assert.NotNil(t, request.Params.Search)
		name, ok := request.Params.Search["name"]
		assert.True(t, ok)
		assert.Equal(t, "Linux*", name)
		assert.True(t, request.Params.SearchWildcardsEnabled)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{GroupID: "1", Name: "Linux servers"},
				{GroupID: "2", Name: "Linux workstations"},
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetSearch(map[string]any{"name": "Linux*"}),
		WithHostGroupGetSearchWildcardsEnabled(true),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Result))
}

// TestHostGroupGetWithGroupIDs tests filtering by group IDs
func TestHostGroupGetWithGroupIDs(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify GroupIDs parameter
		assert.Equal(t, []string{"1", "2", "3"}, request.Params.GroupIDs)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{GroupID: "1", Name: "Group 1"},
				{GroupID: "2", Name: "Group 2"},
				{GroupID: "3", Name: "Group 3"},
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetGroupIDs([]string{"1", "2", "3"}),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Result))
}

// TestHostGroupGetAllFlags tests all boolean flags
func TestHostGroupGetAllFlags(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify all flags are set
		assert.True(t, request.Params.MonitoredHosts)
		assert.True(t, request.Params.RealHosts)
		assert.True(t, request.Params.WithItems)
		assert.True(t, request.Params.WithTriggers)
		assert.True(t, request.Params.Editable)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result:  []HostGroup{},
			ID:      request.ID,
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetMonitoredHosts(true),
		WithHostGroupGetRealHosts(true),
		WithHostGroupGetWithItems(true),
		WithHostGroupGetWithTriggers(true),
		WithHostGroupGetEditable(true),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	_, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
}

// TestHostGroupGetResponseUnmarshaling tests correct response unmarshaling
func TestHostGroupGetResponseUnmarshaling(t *testing.T) {
	t.Parallel()

	// Test that the response unmarshals correctly from JSON
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"groupid": "10",
				"name": "Test Group",
				"flags": "0",
				"internal": "0",
				"uuid": "abc123def456"
			}
		],
		"id": 1
	}`

	var response HostGroupGetResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	assert.NoError(t, err)

	assert.Equal(t, JSONRPC, response.JSONRPC)
	assert.Equal(t, 1, response.ID)
	assert.Nil(t, response.Error)
	assert.Equal(t, 1, len(response.Result))

	// Verify all fields are strings as expected
	group := response.Result[0]
	assert.Equal(t, "10", group.GroupID)
	assert.Equal(t, "Test Group", group.Name)
	assert.Equal(t, "0", group.Flags)
	assert.Equal(t, "0", group.Internal)
	assert.Equal(t, "abc123def456", group.UUID)
}

// TestHostGroupGetAllSelectOptions tests all select options
func TestHostGroupGetAllSelectOptions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify all Select parameters are set
		assert.Equal(t, "extend", request.Params.SelectGroupDiscoveries)
		assert.Equal(t, "extend", request.Params.SelectHosts)
		assert.Equal(t, "extend", request.Params.SelectItems)
		assert.Equal(t, "extend", request.Params.SelectTriggers)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result:  []HostGroup{},
			ID:      request.ID,
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

	request := NewHostGroupGetRequest(
		WithHostGroupGetSelectGroupDiscoveries("extend"),
		WithHostGroupGetSelectHosts("extend"),
		WithHostGroupGetSelectItems("extend"),
		WithHostGroupGetSelectTriggers("extend"),
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	_, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
}

// TestHostGroupGetNewGetAllHostGroupsRequest tests the convenience constructor
func TestHostGroupGetNewGetAllHostGroupsRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request HostGroupGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify output is set to "extend" by default
		assert.Equal(t, "extend", request.Params.Output)

		response := HostGroupGetResponse{
			JSONRPC: JSONRPC,
			Result: []HostGroup{
				{GroupID: "1", Name: "Group 1"},
				{GroupID: "2", Name: "Group 2"},
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

	// Use the convenience constructor
	request := NewGetAllHostGroupsRequest(
		WithHostGroupGetAuth("test-auth-token"),
	)

	ctx := context.Background()
	response, err := client.HostGroupGet(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Result))
}

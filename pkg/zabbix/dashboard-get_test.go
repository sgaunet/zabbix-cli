package zabbix

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDashboardGetRequest(t *testing.T) {
	req := NewDashboardGetRequest()
	if req.JSONRPC != JSONRPC {
		t.Errorf("Expected JSONRPC to be %s, got %s", JSONRPC, req.JSONRPC)
	}
	if req.Method != "dashboard.get" {
		t.Errorf("Expected Method to be dashboard.get, got %s", req.Method)
	}
}

func TestWithDashboardGetDashboardIDs(t *testing.T) {
	ids := []string{"1", "2", "3"}
	req := NewDashboardGetRequest(WithDashboardGetDashboardIDs(ids))
	if len(req.Params.DashboardIDs) != 3 {
		t.Errorf("Expected 3 dashboard IDs, got %d", len(req.Params.DashboardIDs))
	}
	if req.Params.DashboardIDs[0] != "1" {
		t.Errorf("Expected first dashboard ID to be 1, got %s", req.Params.DashboardIDs[0])
	}
}

func TestWithDashboardGetSelectPages(t *testing.T) {
	req := NewDashboardGetRequest(WithDashboardGetSelectPages("extend"))
	if req.Params.SelectPages != "extend" {
		t.Errorf("Expected SelectPages to be extend, got %v", req.Params.SelectPages)
	}
}

func TestWithDashboardGetOutput(t *testing.T) {
	req := NewDashboardGetRequest(WithDashboardGetOutput("extend"))
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
}

func TestWithDashboardGetAuth(t *testing.T) {
	token := "test-auth-token"
	req := NewDashboardGetRequest(WithDashboardGetAuth(token))
	if req.Auth != token {
		t.Errorf("Expected Auth to be %s, got %s", token, req.Auth)
	}
}

func TestWithDashboardGetID(t *testing.T) {
	id := 123
	req := NewDashboardGetRequest(WithDashboardGetID(id))
	if req.ID != id {
		t.Errorf("Expected ID to be %d, got %d", id, req.ID)
	}
}

func TestWithDashboardGetFilter(t *testing.T) {
	filter := map[string]interface{}{
		"name": "test dashboard",
	}
	req := NewDashboardGetRequest(WithDashboardGetFilter(filter))
	if req.Params.Filter["name"] != "test dashboard" {
		t.Errorf("Expected filter name to be 'test dashboard', got %v", req.Params.Filter["name"])
	}
}

func TestNewGetAllDashboardsRequest(t *testing.T) {
	req := NewGetAllDashboardsRequest()
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
	if req.Params.SelectPages != "extend" {
		t.Errorf("Expected SelectPages to be extend, got %v", req.Params.SelectPages)
	}
}

func TestNewGetAllDashboardsRequestWithOverrides(t *testing.T) {
	token := "override-token"
	id := 999
	req := NewGetAllDashboardsRequest(
		WithDashboardGetAuth(token),
		WithDashboardGetID(id),
	)
	if req.Auth != token {
		t.Errorf("Expected Auth to be %s, got %s", token, req.Auth)
	}
	if req.ID != id {
		t.Errorf("Expected ID to be %d, got %d", id, req.ID)
	}
	if req.Params.Output != "extend" {
		t.Errorf("Expected Output to be extend, got %v", req.Params.Output)
	}
}

// TestWithDashboardGetSelectPagesArray tests SelectPages with array of field names
func TestWithDashboardGetSelectPagesArray(t *testing.T) {
	t.Parallel()

	fields := []string{"dashboard_pageid", "name", "display_period"}
	req := NewDashboardGetRequest(WithDashboardGetSelectPages(fields))

	assert.Equal(t, fields, req.Params.SelectPages)
}

// TestWithDashboardGetSelectUsers tests SelectUsers with "extend"
func TestWithDashboardGetSelectUsers(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetSelectUsers("extend"))
	assert.Equal(t, "extend", req.Params.SelectUsers)
}

// TestWithDashboardGetSelectUsersArray tests SelectUsers with array of fields
func TestWithDashboardGetSelectUsersArray(t *testing.T) {
	t.Parallel()

	fields := []string{"userid", "permission"}
	req := NewDashboardGetRequest(WithDashboardGetSelectUsers(fields))

	assert.Equal(t, fields, req.Params.SelectUsers)
}

// TestWithDashboardGetSelectUserGroups tests SelectUserGroups with "extend"
func TestWithDashboardGetSelectUserGroups(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetSelectUserGroups("extend"))
	assert.Equal(t, "extend", req.Params.SelectUserGroups)
}

// TestWithDashboardGetSelectUserGroupsArray tests SelectUserGroups with array of fields
func TestWithDashboardGetSelectUserGroupsArray(t *testing.T) {
	t.Parallel()

	fields := []string{"usrgrpid", "permission"}
	req := NewDashboardGetRequest(WithDashboardGetSelectUserGroups(fields))

	assert.Equal(t, fields, req.Params.SelectUserGroups)
}

// TestDashboardGetAllSelectParameters tests all select parameters combined
func TestDashboardGetAllSelectParameters(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(
		WithDashboardGetSelectPages("extend"),
		WithDashboardGetSelectUsers([]string{"userid", "permission"}),
		WithDashboardGetSelectUserGroups("extend"),
	)

	assert.Equal(t, "extend", req.Params.SelectPages)
	assert.Equal(t, []string{"userid", "permission"}, req.Params.SelectUsers)
	assert.Equal(t, "extend", req.Params.SelectUserGroups)
}

// TestDashboardGetJSONMarshalingWithSelectParams validates JSON marshaling of select parameters
func TestDashboardGetJSONMarshalingWithSelectParams(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(
		WithDashboardGetOutput("extend"),
		WithDashboardGetSelectPages([]string{"dashboard_pageid", "name"}),
		WithDashboardGetSelectUsers("extend"),
		WithDashboardGetSelectUserGroups([]string{"usrgrpid"}),
	)

	data, err := json.Marshal(req.Params)
	assert.NoError(t, err)

	var params map[string]any
	err = json.Unmarshal(data, &params)
	assert.NoError(t, err)

	// Verify select parameters are present in JSON
	assert.Equal(t, "extend", params["output"])

	// SelectPages should be an array
	selectPages, ok := params["selectPages"].([]any)
	assert.True(t, ok, "selectPages should be an array")
	assert.Equal(t, 2, len(selectPages))
	assert.Equal(t, "dashboard_pageid", selectPages[0])
	assert.Equal(t, "name", selectPages[1])

	// SelectUsers should be "extend"
	assert.Equal(t, "extend", params["selectUsers"])

	// SelectUserGroups should be an array
	selectUserGroups, ok := params["selectUserGroups"].([]any)
	assert.True(t, ok, "selectUserGroups should be an array")
	assert.Equal(t, 1, len(selectUserGroups))
	assert.Equal(t, "usrgrpid", selectUserGroups[0])
}

// TestDashboardGetWithLimit tests the Limit parameter
func TestDashboardGetWithLimit(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetLimit(10))
	assert.Equal(t, 10, req.Params.Limit)
}

// TestDashboardGetWithSortField tests sorting parameters
func TestDashboardGetWithSortField(t *testing.T) {
	t.Parallel()

	sortField := []string{"name", "dashboardid"}
	sortOrder := []string{"ASC", "DESC"}

	req := NewDashboardGetRequest(
		WithDashboardGetSortField(sortField),
		WithDashboardGetSortOrder(sortOrder),
	)

	assert.Equal(t, sortField, req.Params.SortField)
	assert.Equal(t, sortOrder, req.Params.SortOrder)
}

// TestDashboardGetWithSearch tests search parameters
func TestDashboardGetWithSearch(t *testing.T) {
	t.Parallel()

	search := map[string]any{"name": "Dashboard"}
	req := NewDashboardGetRequest(
		WithDashboardGetSearch(search),
		WithDashboardGetSearchWildcardsEnabled(true),
	)

	assert.Equal(t, "Dashboard", req.Params.Search["name"])
	assert.True(t, req.Params.SearchWildcardsEnabled)
}

// TestDashboardGetWithSearchByAny tests searchByAny parameter
func TestDashboardGetWithSearchByAny(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetSearchByAny(true))
	assert.True(t, req.Params.SearchByAny)
}

// TestDashboardGetWithStartSearch tests startSearch parameter
func TestDashboardGetWithStartSearch(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetStartSearch(true))
	assert.True(t, req.Params.StartSearch)
}

// TestDashboardGetWithExcludeSearch tests excludeSearch parameter
func TestDashboardGetWithExcludeSearch(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetExcludeSearch(true))
	assert.True(t, req.Params.ExcludeSearch)
}

// TestDashboardGetWithPreserveKeys tests preservekeys parameter
func TestDashboardGetWithPreserveKeys(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetPreserveKeys(true))
	assert.True(t, req.Params.PreserveKeys)
}

// TestDashboardGetWithCountOutput tests countOutput parameter
func TestDashboardGetWithCountOutput(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(WithDashboardGetCountOutput(true))
	assert.True(t, req.Params.CountOutput)
}

// TestDashboardGetCompleteRequest tests a complete request with all parameters
func TestDashboardGetCompleteRequest(t *testing.T) {
	t.Parallel()

	dashboardIDs := []string{"1", "2", "3"}
	filter := map[string]any{"private": "0"}

	req := NewDashboardGetRequest(
		WithDashboardGetDashboardIDs(dashboardIDs),
		WithDashboardGetOutput("extend"),
		WithDashboardGetSelectPages("extend"),
		WithDashboardGetSelectUsers([]string{"userid"}),
		WithDashboardGetSelectUserGroups("extend"),
		WithDashboardGetFilter(filter),
		WithDashboardGetLimit(50),
		WithDashboardGetSortField([]string{"name"}),
		WithDashboardGetSortOrder([]string{"ASC"}),
		WithDashboardGetAuth("test-token"),
		WithDashboardGetID(1),
	)

	// Verify all parameters are set correctly
	assert.Equal(t, dashboardIDs, req.Params.DashboardIDs)
	assert.Equal(t, "extend", req.Params.Output)
	assert.Equal(t, "extend", req.Params.SelectPages)
	assert.Equal(t, []string{"userid"}, req.Params.SelectUsers)
	assert.Equal(t, "extend", req.Params.SelectUserGroups)
	assert.Equal(t, "0", req.Params.Filter["private"])
	assert.Equal(t, 50, req.Params.Limit)
	assert.Equal(t, []string{"name"}, req.Params.SortField)
	assert.Equal(t, []string{"ASC"}, req.Params.SortOrder)
	assert.Equal(t, "test-token", req.Auth)
	assert.Equal(t, 1, req.ID)
}

// TestDashboardGetRequestJSONStructure validates complete JSON structure
func TestDashboardGetRequestJSONStructure(t *testing.T) {
	t.Parallel()

	req := NewDashboardGetRequest(
		WithDashboardGetDashboardIDs([]string{"100"}),
		WithDashboardGetOutput("extend"),
		WithDashboardGetSelectPages("extend"),
		WithDashboardGetAuth("auth-token"),
		WithDashboardGetID(5),
	)

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var jsonReq map[string]any
	err = json.Unmarshal(data, &jsonReq)
	assert.NoError(t, err)

	// Verify JSON-RPC structure
	assert.Equal(t, JSONRPC, jsonReq["jsonrpc"])
	assert.Equal(t, "dashboard.get", jsonReq["method"])
	assert.Equal(t, "auth-token", jsonReq["auth"])
	assert.Equal(t, float64(5), jsonReq["id"])

	// Verify params structure
	params, ok := jsonReq["params"].(map[string]any)
	assert.True(t, ok, "params should be a map")

	dashboardIDs, ok := params["dashboardids"].([]any)
	assert.True(t, ok, "dashboardids should be an array")
	assert.Equal(t, 1, len(dashboardIDs))
	assert.Equal(t, "100", dashboardIDs[0])

	assert.Equal(t, "extend", params["output"])
	assert.Equal(t, "extend", params["selectPages"])
}

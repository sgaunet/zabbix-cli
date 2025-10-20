package zabbix

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDashboardPageStructure validates the DashboardPage structure matches API spec
func TestDashboardPageStructure(t *testing.T) {
	t.Parallel()

	page := DashboardPage{
		DashboardPageID: "123",
		Name:            "Page 1",
		DisplayPeriod:   "60",
		Widgets: []Widget{
			{
				WidgetID: "456",
				Type:     "problems",
				Name:     "Current problems",
				X:        "0",
				Y:        "0",
				Width:    "12",
				Height:   "5",
				View:     "0",
				Fields: []WidgetField{
					{Type: "0", Name: "severities", Value: "4"},
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(page)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled DashboardPage
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Verify structure
	assert.Equal(t, "123", unmarshaled.DashboardPageID)
	assert.Equal(t, "Page 1", unmarshaled.Name)
	assert.Equal(t, "60", unmarshaled.DisplayPeriod)
	assert.Equal(t, 1, len(unmarshaled.Widgets))
	assert.Equal(t, "456", unmarshaled.Widgets[0].WidgetID)
	assert.Equal(t, "problems", unmarshaled.Widgets[0].Type)
}

// TestDashboardPageJSONTags validates correct JSON field names
func TestDashboardPageJSONTags(t *testing.T) {
	t.Parallel()

	page := DashboardPage{
		DashboardPageID: "1",
		Name:            "Main Page",
		DisplayPeriod:   "30",
		Widgets:         []Widget{},
	}

	data, err := json.Marshal(page)
	assert.NoError(t, err)

	// Parse as map to check JSON keys
	var jsonMap map[string]any
	err = json.Unmarshal(data, &jsonMap)
	assert.NoError(t, err)

	// Verify correct JSON field names according to API spec
	assert.Equal(t, "1", jsonMap["dashboard_pageid"])
	assert.Equal(t, "Main Page", jsonMap["name"])
	assert.Equal(t, "30", jsonMap["display_period"])
	// Empty widgets array should be omitted due to omitempty tag
	_, hasWidgets := jsonMap["widgets"]
	assert.False(t, hasWidgets, "widgets field should be omitted when empty")
}

// TestDashboardWithMultiplePages validates dashboard with multiple pages
func TestDashboardWithMultiplePages(t *testing.T) {
	t.Parallel()

	dashboard := Dashboard{
		DashboardID:   "100",
		Name:          "Multi-page Dashboard",
		DisplayPeriod: "60",
		Pages: []DashboardPage{
			{
				DashboardPageID: "1",
				Name:            "Overview",
				DisplayPeriod:   "30",
				Widgets: []Widget{
					{WidgetID: "10", Type: "graph", Name: "CPU Usage"},
					{WidgetID: "11", Type: "graph", Name: "Memory Usage"},
				},
			},
			{
				DashboardPageID: "2",
				Name:            "Details",
				DisplayPeriod:   "60",
				Widgets: []Widget{
					{WidgetID: "20", Type: "problems", Name: "Critical Issues"},
				},
			},
		},
	}

	// Marshal and unmarshal
	data, err := json.Marshal(dashboard)
	assert.NoError(t, err)

	var unmarshaled Dashboard
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Validate structure
	assert.Equal(t, "100", unmarshaled.DashboardID)
	assert.Equal(t, 2, len(unmarshaled.Pages))
	assert.Equal(t, "Overview", unmarshaled.Pages[0].Name)
	assert.Equal(t, 2, len(unmarshaled.Pages[0].Widgets))
	assert.Equal(t, "Details", unmarshaled.Pages[1].Name)
	assert.Equal(t, 1, len(unmarshaled.Pages[1].Widgets))
}

// TestWidgetStructure validates widget nesting and field structure
func TestWidgetStructure(t *testing.T) {
	t.Parallel()

	widget := Widget{
		WidgetID: "789",
		Type:     "problems",
		Name:     "Active Problems",
		X:        "0",
		Y:        "0",
		Width:    "24",
		Height:   "10",
		View:     "0",
		Fields: []WidgetField{
			{Type: "2", Name: "groupids", Value: "4"},
			{Type: "3", Name: "hostids", Value: "10084"},
			{Type: "0", Name: "severities", Value: "4"},
			{Type: "0", Name: "severities", Value: "5"},
		},
	}

	data, err := json.Marshal(widget)
	assert.NoError(t, err)

	var unmarshaled Widget
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, "789", unmarshaled.WidgetID)
	assert.Equal(t, "problems", unmarshaled.Type)
	assert.Equal(t, 4, len(unmarshaled.Fields))
	assert.Equal(t, "groupids", unmarshaled.Fields[0].Name)
}

// TestWidgetJSONTags validates widget JSON field names
func TestWidgetJSONTags(t *testing.T) {
	t.Parallel()

	widget := Widget{
		WidgetID: "1",
		Type:     "graph",
		Name:     "Test Widget",
		X:        "0",
		Y:        "0",
		Width:    "12",
		Height:   "5",
		View:     "0",
		Fields:   []WidgetField{},
	}

	data, err := json.Marshal(widget)
	assert.NoError(t, err)

	var jsonMap map[string]any
	err = json.Unmarshal(data, &jsonMap)
	assert.NoError(t, err)

	// Verify JSON field names
	assert.Equal(t, "1", jsonMap["widgetid"])
	assert.Equal(t, "graph", jsonMap["type"])
	assert.Equal(t, "Test Widget", jsonMap["name"])
	assert.Equal(t, "0", jsonMap["x"])
	assert.Equal(t, "0", jsonMap["y"])
	assert.Equal(t, "12", jsonMap["width"])
	assert.Equal(t, "5", jsonMap["height"])
	assert.Equal(t, "0", jsonMap["view_mode"])
	// Empty fields array should be omitted due to omitempty tag
	_, hasFields := jsonMap["fields"]
	assert.False(t, hasFields, "fields should be omitted when empty")
}

// TestWidgetFieldTypes validates widget field value types
func TestWidgetFieldTypes(t *testing.T) {
	t.Parallel()

	fields := []WidgetField{
		{Type: "0", Name: "integer_field", Value: 42},
		{Type: "1", Name: "string_field", Value: "test"},
		{Type: "2", Name: "array_field", Value: []string{"1", "2", "3"}},
	}

	data, err := json.Marshal(fields)
	assert.NoError(t, err)

	var unmarshaled []WidgetField
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(unmarshaled))
	// Integer value
	assert.Equal(t, float64(42), unmarshaled[0].Value) // JSON unmarshals numbers as float64
	// String value
	assert.Equal(t, "test", unmarshaled[1].Value)
	// Array value (unmarshals as []interface{})
	arrayVal, ok := unmarshaled[2].Value.([]any)
	assert.True(t, ok, "array value should unmarshal as []any")
	assert.Equal(t, 3, len(arrayVal))
}

// TestDashboardGetResponseWithPages validates response unmarshaling with pages
func TestDashboardGetResponseWithPages(t *testing.T) {
	t.Parallel()

	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"dashboardid": "1",
				"name": "Test Dashboard",
				"display_period": "30",
				"pages": [
					{
						"dashboard_pageid": "10",
						"name": "Page 1",
						"display_period": "60",
						"widgets": [
							{
								"widgetid": "100",
								"type": "problems",
								"name": "Current problems",
								"x": "0",
								"y": "0",
								"width": "12",
								"height": "5",
								"view_mode": "0",
								"fields": [
									{
										"type": "0",
										"name": "severities",
										"value": "4"
									}
								]
							}
						]
					}
				]
			}
		],
		"id": 1
	}`

	var response DashboardGetResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	assert.NoError(t, err)

	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Equal(t, 1, len(response.Result))

	dashboard := response.Result[0]
	assert.Equal(t, "1", dashboard.DashboardID)
	assert.Equal(t, "Test Dashboard", dashboard.Name)
	assert.Equal(t, 1, len(dashboard.Pages))

	page := dashboard.Pages[0]
	assert.Equal(t, "10", page.DashboardPageID)
	assert.Equal(t, "Page 1", page.Name)
	assert.Equal(t, 1, len(page.Widgets))

	widget := page.Widgets[0]
	assert.Equal(t, "100", widget.WidgetID)
	assert.Equal(t, "problems", widget.Type)
	assert.Equal(t, 1, len(widget.Fields))

	field := widget.Fields[0]
	assert.Equal(t, "0", field.Type)
	assert.Equal(t, "severities", field.Name)
	assert.Equal(t, "4", field.Value)
}

// TestEmptyDashboardPages validates handling of dashboards without pages
func TestEmptyDashboardPages(t *testing.T) {
	t.Parallel()

	dashboard := Dashboard{
		DashboardID: "1",
		Name:        "Empty Dashboard",
		Pages:       []DashboardPage{},
	}

	data, err := json.Marshal(dashboard)
	assert.NoError(t, err)

	var unmarshaled Dashboard
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, "1", unmarshaled.DashboardID)
	assert.Equal(t, 0, len(unmarshaled.Pages))
}

// TestWidgetPositioning validates widget X/Y positioning
func TestWidgetPositioning(t *testing.T) {
	t.Parallel()

	page := DashboardPage{
		DashboardPageID: "1",
		Widgets: []Widget{
			{WidgetID: "1", Type: "graph", X: "0", Y: "0", Width: "12", Height: "5"},
			{WidgetID: "2", Type: "graph", X: "12", Y: "0", Width: "12", Height: "5"},
			{WidgetID: "3", Type: "problems", X: "0", Y: "5", Width: "24", Height: "10"},
		},
	}

	data, err := json.Marshal(page)
	assert.NoError(t, err)

	var unmarshaled DashboardPage
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	// Verify widget positions
	assert.Equal(t, "0", unmarshaled.Widgets[0].X)
	assert.Equal(t, "0", unmarshaled.Widgets[0].Y)
	assert.Equal(t, "12", unmarshaled.Widgets[1].X)
	assert.Equal(t, "0", unmarshaled.Widgets[1].Y)
	assert.Equal(t, "0", unmarshaled.Widgets[2].X)
	assert.Equal(t, "5", unmarshaled.Widgets[2].Y)
}

// TestOmitEmptyFields validates omitempty behavior
func TestOmitEmptyFields(t *testing.T) {
	t.Parallel()

	// Minimal widget (only required fields)
	widget := Widget{
		Type: "problems",
	}

	data, err := json.Marshal(widget)
	assert.NoError(t, err)

	var jsonMap map[string]any
	err = json.Unmarshal(data, &jsonMap)
	assert.NoError(t, err)

	// Type is required, should be present
	assert.Equal(t, "problems", jsonMap["type"])

	// Optional fields with omitempty should not be present if empty
	_, hasWidgetID := jsonMap["widgetid"]
	assert.False(t, hasWidgetID, "widgetid should be omitted when empty")

	_, hasName := jsonMap["name"]
	assert.False(t, hasName, "name should be omitted when empty")
}

// TestDashboardUserGroups validates dashboard sharing structures
func TestDashboardUserGroups(t *testing.T) {
	t.Parallel()

	dashboard := Dashboard{
		DashboardID: "1",
		Name:        "Shared Dashboard",
		Users: []DashboardUser{
			{UserID: "10", Permission: "2"},
			{UserID: "20", Permission: "3"},
		},
		UserGroups: []DashboardUserGroup{
			{UserGroupID: "5", Permission: "2"},
		},
	}

	data, err := json.Marshal(dashboard)
	assert.NoError(t, err)

	var unmarshaled Dashboard
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(unmarshaled.Users))
	assert.Equal(t, "10", unmarshaled.Users[0].UserID)
	assert.Equal(t, "2", unmarshaled.Users[0].Permission)

	assert.Equal(t, 1, len(unmarshaled.UserGroups))
	assert.Equal(t, "5", unmarshaled.UserGroups[0].UserGroupID)
}

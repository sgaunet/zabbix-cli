package zabbix_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProblemGetResponseIntegration tests Problem parsing with realistic API 7.2 response data
func TestProblemGetResponseIntegration(t *testing.T) {
	// Realistic Zabbix API 7.2 problem.get response with mixed types
	responseJSON := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"eventid": "12345",
				"source": "0",
				"object": "0",
				"objectid": "67890",
				"clock": "1609459200",
				"ns": "123456789",
				"r_eventid": "12346",
				"r_clock": "1609462800",
				"r_ns": "987654321",
				"name": "High CPU usage on server-01",
				"acknowledged": "1",
				"severity": "4",
				"suppressed": "0",
				"opdata": "CPU: 95%"
			},
			{
				"eventid": "12347",
				"source": "0",
				"object": "0",
				"objectid": "67891",
				"clock": "1609460000",
				"ns": "111222333",
				"r_eventid": "0",
				"r_clock": "0",
				"r_ns": "0",
				"name": "Disk space low on server-02",
				"acknowledged": "0",
				"severity": "3",
				"suppressed": "0",
				"opdata": "Free: 5%"
			}
		],
		"id": 1
	}`

	// Define response structure inline since GetProblemResponse may not be exported
	var response struct {
		JSONRPC string           `json:"jsonrpc"`
		Result  []zabbix.Problem `json:"result"`
		ID      int              `json:"id"`
	}
	err := json.Unmarshal([]byte(responseJSON), &response)
	require.NoError(t, err, "Failed to unmarshal problem response")

	// Verify we got 2 problems
	assert.Len(t, response.Result, 2, "Expected 2 problems in response")

	// Test first problem (acknowledged)
	problem1 := response.Result[0]
	assert.Equal(t, "12345", problem1.EventID, "EventID should be string")
	assert.Equal(t, "67890", problem1.ObjectID, "ObjectID should be string")
	assert.True(t, problem1.GetAcknowledge(), "Problem should be acknowledged")
	assert.False(t, problem1.GetSuppressed(), "Problem should not be suppressed")
	assert.Equal(t, "High", problem1.GetSeverity(), "Severity 4 should be High")
	assert.Equal(t, "Yes", problem1.GetAcknowledgeStr(), "Acknowledged should be Yes")

	// Test timestamp conversion
	clockTime := problem1.GetClock()
	assert.False(t, clockTime.IsZero(), "Clock time should not be zero")
	assert.Equal(t, int64(1609459200), problem1.Clock.Int64(), "Clock should be 1609459200")

	rClockTime := problem1.GetRClock()
	assert.False(t, rClockTime.IsZero(), "R_Clock time should not be zero")
	assert.Equal(t, int64(1609462800), problem1.Rclock.Int64(), "R_Clock should be 1609462800")

	// Test duration calculation
	duration := problem1.GetDuration()
	expectedDuration := time.Duration(3600) * time.Second // 1 hour difference
	assert.Equal(t, expectedDuration, duration, "Duration should be 1 hour")

	// Test second problem (not acknowledged, not resolved)
	problem2 := response.Result[1]
	assert.Equal(t, "12347", problem2.EventID, "EventID should be string")
	assert.False(t, problem2.GetAcknowledge(), "Problem should not be acknowledged")
	assert.Equal(t, "Average", problem2.GetSeverity(), "Severity 3 should be Average")
	assert.Equal(t, "No", problem2.GetAcknowledgeStr(), "Acknowledged should be No")

	// Verify R_Clock is zero (not resolved)
	rClock2 := problem2.GetRClock()
	assert.True(t, rClock2.IsZero(), "R_Clock should be zero for unresolved problem")
}

// TestMaintenanceGetResponseIntegration tests Maintenance parsing with realistic API 7.2 response
func TestMaintenanceGetResponseIntegration(t *testing.T) {
	// Realistic Zabbix API 7.2 maintenance.get response
	responseJSON := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"maintenanceid": "100",
				"name": "Weekend Maintenance",
				"active_since": "1609459200",
				"active_till": "1609545600",
				"description": "Regular weekend maintenance window",
				"maintenance_type": "1",
				"created_at": "1609300000",
				"updated_at": "1609400000",
				"timeperiods": [
					{
						"timeperiodid": "200",
						"timeperiod_type": 2,
						"start_date": 0,
						"period": 7200,
						"dayofweek": 0,
						"start_time": 28800
					}
				],
				"groupids": ["10", "20"],
				"hostids": ["30", "40", "50"],
				"tags": [
					{
						"tag": "Environment",
						"value": "Production",
						"operator": 0
					}
				],
				"tags_evaltype": 0
			}
		],
		"id": 1
	}`

	var response struct {
		JSONRPC string                `json:"jsonrpc"`
		Result  []zabbix.Maintenance  `json:"result"`
		ID      int                   `json:"id"`
		Error   *zabbix.Error         `json:"error,omitempty"`
	}

	err := json.Unmarshal([]byte(responseJSON), &response)
	require.NoError(t, err, "Failed to unmarshal maintenance response")

	// Verify maintenance data
	assert.Len(t, response.Result, 1, "Expected 1 maintenance in response")

	maint := response.Result[0]

	// Verify ID fields are strings
	assert.Equal(t, "100", maint.MaintenanceID, "MaintenanceID should be string")
	assert.Equal(t, "Weekend Maintenance", maint.Name)

	// Verify timestamps
	assert.Equal(t, int64(1609459200), maint.ActiveSince.Int64(), "ActiveSince timestamp")
	assert.Equal(t, int64(1609545600), maint.ActiveTill.Int64(), "ActiveTill timestamp")
	assert.Equal(t, int64(1609300000), maint.CreatedAt.Int64(), "CreatedAt timestamp")
	assert.Equal(t, int64(1609400000), maint.UpdatedAt.Int64(), "UpdatedAt timestamp")

	// Verify maintenance type
	assert.Equal(t, zabbix.MaintenanceNoDataCollection, maint.MaintenanceType)

	// Verify time periods
	require.Len(t, maint.TimePeriods, 1, "Expected 1 time period")
	tp := maint.TimePeriods[0]
	assert.Equal(t, "200", tp.TimePeriodID, "TimePeriodID should be string")
	assert.Equal(t, zabbix.TimePeriodTypeWeekly, tp.TimePeriodType)
	assert.Equal(t, 7200, tp.Period, "Period should be 7200 seconds (2 hours)")

	// Verify groupids and hostids are string arrays
	assert.Equal(t, []string{"10", "20"}, maint.GroupIDs, "GroupIDs should be string array")
	assert.Equal(t, []string{"30", "40", "50"}, maint.HostIDs, "HostIDs should be string array")

	// Verify tags
	require.Len(t, maint.Tags, 1, "Expected 1 tag")
	assert.Equal(t, "Environment", maint.Tags[0].Tag)
	assert.Equal(t, "Production", maint.Tags[0].Value)
}

// TestDashboardGetResponseIntegration tests Dashboard parsing with realistic API 7.2 response including widgets
func TestDashboardGetResponseIntegration(t *testing.T) {
	// Realistic Zabbix API 7.2 dashboard.get response with complex nested widgets
	responseJSON := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"dashboardid": "500",
				"name": "Infrastructure Overview",
				"userid": "1",
				"private": "0",
				"display_period": "60",
				"auto_start": "1",
				"pages": [
					{
						"dashboard_pageid": "600",
						"name": "Main View",
						"display_period": "0",
						"widgets": [
							{
								"widgetid": "700",
								"type": "graph",
								"name": "CPU Usage",
								"x": "0",
								"y": "0",
								"width": "12",
								"height": "5",
								"view_mode": "0",
								"fields": [
									{
										"type": "0",
										"name": "graphid",
										"value": "12345"
									},
									{
										"type": "0",
										"name": "reference",
										"value": "ABCDE"
									}
								]
							},
							{
								"widgetid": "701",
								"type": "problems",
								"name": "Active Problems",
								"x": "12",
								"y": "0",
								"width": "12",
								"height": "5",
								"view_mode": "0",
								"fields": [
									{
										"type": "3",
										"name": "show_tags",
										"value": "3"
									},
									{
										"type": "0",
										"name": "severities",
										"value": "4,5"
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

	var response zabbix.DashboardGetResponse
	err := json.Unmarshal([]byte(responseJSON), &response)
	require.NoError(t, err, "Failed to unmarshal dashboard response")

	// Verify dashboard data
	assert.Len(t, response.Result, 1, "Expected 1 dashboard in response")

	dashboard := response.Result[0]

	// Verify ID fields are strings
	assert.Equal(t, "500", dashboard.DashboardID, "DashboardID should be string")
	assert.Equal(t, "1", dashboard.UserID, "UserID should be string")
	assert.Equal(t, "Infrastructure Overview", dashboard.Name)

	// Verify boolean fields using BoolString
	assert.False(t, dashboard.Private.Bool(), "Private should be false (\"0\")")
	assert.True(t, dashboard.AutoStart.Bool(), "AutoStart should be true (\"1\")")

	// Verify dashboard display period
	assert.Equal(t, "60", dashboard.DisplayPeriod, "DisplayPeriod should be string \"60\"")

	// Verify pages
	require.Len(t, dashboard.Pages, 1, "Expected 1 page")
	page := dashboard.Pages[0]
	assert.Equal(t, "600", page.DashboardPageID, "DashboardPageID should be string")
	assert.Equal(t, "Main View", page.Name)

	// Verify widgets
	require.Len(t, page.Widgets, 2, "Expected 2 widgets")

	// Verify first widget (graph)
	widget1 := page.Widgets[0]
	assert.Equal(t, "700", widget1.WidgetID, "WidgetID should be string")
	assert.Equal(t, "graph", widget1.Type)
	assert.Equal(t, "CPU Usage", widget1.Name)
	assert.Equal(t, "0", widget1.X, "X coordinate should be string")
	assert.Equal(t, "0", widget1.Y, "Y coordinate should be string")
	assert.Equal(t, "12", widget1.Width, "Width should be string")
	assert.Equal(t, "5", widget1.Height, "Height should be string")

	// Verify widget fields
	require.Len(t, widget1.Fields, 2, "Expected 2 fields in first widget")
	assert.Equal(t, "0", widget1.Fields[0].Type, "Field type should be string")
	assert.Equal(t, "graphid", widget1.Fields[0].Name)
	assert.Equal(t, "12345", widget1.Fields[0].Value)

	// Verify second widget (problems)
	widget2 := page.Widgets[1]
	assert.Equal(t, "701", widget2.WidgetID, "WidgetID should be string")
	assert.Equal(t, "problems", widget2.Type)
	assert.Equal(t, "Active Problems", widget2.Name)
	assert.Len(t, widget2.Fields, 2, "Expected 2 fields in second widget")
}

// TestHostGroupGetResponseIntegration tests HostGroup parsing with realistic API 7.2 response
func TestHostGroupGetResponseIntegration(t *testing.T) {
	responseJSON := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"groupid": "10",
				"name": "Linux servers",
				"flags": "0",
				"internal": "0",
				"uuid": "dc84c7bf78274b8e94d0a89cebfc41ff"
			},
			{
				"groupid": "20",
				"name": "Zabbix servers",
				"flags": "0",
				"internal": "1",
				"uuid": "a89c8cc7eb5e419da29a6c25a0c0cc77"
			}
		],
		"id": 1
	}`

	var response zabbix.HostGroupGetResponse
	err := json.Unmarshal([]byte(responseJSON), &response)
	require.NoError(t, err, "Failed to unmarshal hostgroup response")

	// Verify hostgroups
	assert.Len(t, response.Result, 2, "Expected 2 hostgroups in response")

	// Verify first hostgroup
	hg1 := response.Result[0]
	assert.Equal(t, "10", hg1.GroupID, "GroupID should be string")
	assert.Equal(t, "Linux servers", hg1.Name)
	assert.Equal(t, "0", hg1.Flags, "Flags should be string")
	assert.Equal(t, "0", hg1.Internal, "Internal should be string")
	assert.Equal(t, "dc84c7bf78274b8e94d0a89cebfc41ff", hg1.UUID)

	// Verify second hostgroup (internal)
	hg2 := response.Result[1]
	assert.Equal(t, "20", hg2.GroupID, "GroupID should be string")
	assert.Equal(t, "Zabbix servers", hg2.Name)
	assert.Equal(t, "1", hg2.Internal, "Internal flag should be \"1\" for internal group")
}

// TestRoundTripProblemRequest tests Problem request marshaling with all type fixes
func TestRoundTripProblemRequest(t *testing.T) {
	// Create a request with all field types
	request := zabbix.GetProblemRequest{
		JSONRPC: "2.0",
		Method:  "problem.get",
		Params: zabbix.ProblemParams{
			CommonGetParams: zabbix.CommonGetParams{
				Limit:       100,
				CountOutput: false,
			},
			EventIDs:     []string{"12345", "67890"},
			GroupsIDs:    []string{"10", "20"},
			HostsIDs:     []string{"30", "40"},
			Severities:   []string{"3", "4", "5"},
			TimeFrom:     1609459200,
			TimeTill:     1609545600,
			Acknowledged: true,
			Suppressed:   false,
		},
		ID: 1,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(request)
	require.NoError(t, err, "Failed to marshal request")

	// Unmarshal back
	var roundTrip zabbix.GetProblemRequest
	err = json.Unmarshal(jsonData, &roundTrip)
	require.NoError(t, err, "Failed to unmarshal request")

	// Verify all fields survived round-trip
	assert.Equal(t, request.JSONRPC, roundTrip.JSONRPC)
	assert.Equal(t, request.Method, roundTrip.Method)
	assert.Equal(t, request.Params.EventIDs, roundTrip.Params.EventIDs, "EventIDs should survive round-trip")
	assert.Equal(t, request.Params.Severities, roundTrip.Params.Severities, "Severities should survive round-trip")
	assert.Equal(t, request.Params.TimeFrom, roundTrip.Params.TimeFrom, "TimeFrom should survive round-trip")
	assert.Equal(t, request.Params.Acknowledged, roundTrip.Params.Acknowledged, "Acknowledged should survive round-trip")
	assert.Equal(t, request.Params.CommonGetParams.Limit, roundTrip.Params.CommonGetParams.Limit, "Limit should survive round-trip")
}

// TestRoundTripMaintenanceCreate tests Maintenance creation request with timestamp handling
func TestRoundTripMaintenanceCreate(t *testing.T) {
	// Create maintenance with all field types
	now := time.Now().Unix()

	maintenance := zabbix.Maintenance{
		Name:            "Test Maintenance",
		ActiveSince:     zabbix.StringInt64(now),
		ActiveTill:      zabbix.StringInt64(now + 3600),
		Description:     "Test maintenance window",
		MaintenanceType: zabbix.MaintenanceNoDataCollection,
		GroupIDs:        []string{"10", "20"},
		HostIDs:         []string{"30"},
		TimePeriods: []zabbix.TimePeriod{
			{
				TimePeriodType: zabbix.TimePeriodTypeOneTime,
				StartDate:      now,
				Period:         3600,
			},
		},
		Tags: []zabbix.ProblemTag{
			{
				Tag:      "Environment",
				Value:    "Production",
				Operator: 0,
			},
		},
		TagsEvalType: 0,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(maintenance)
	require.NoError(t, err, "Failed to marshal maintenance")

	// Unmarshal back
	var roundTrip zabbix.Maintenance
	err = json.Unmarshal(jsonData, &roundTrip)
	require.NoError(t, err, "Failed to unmarshal maintenance")

	// Verify all fields survived round-trip
	assert.Equal(t, maintenance.Name, roundTrip.Name)
	assert.Equal(t, maintenance.ActiveSince.Int64(), roundTrip.ActiveSince.Int64(), "ActiveSince timestamp should survive")
	assert.Equal(t, maintenance.ActiveTill.Int64(), roundTrip.ActiveTill.Int64(), "ActiveTill timestamp should survive")
	assert.Equal(t, maintenance.MaintenanceType, roundTrip.MaintenanceType)
	assert.Equal(t, maintenance.GroupIDs, roundTrip.GroupIDs, "GroupIDs array should survive")
	assert.Equal(t, maintenance.HostIDs, roundTrip.HostIDs, "HostIDs array should survive")
	assert.Len(t, roundTrip.TimePeriods, 1, "TimePeriods should survive")
	assert.Len(t, roundTrip.Tags, 1, "Tags should survive")
}

// TestStringInt64EdgeCases tests StringInt64 handling of edge cases
func TestStringInt64EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected int64
	}{
		{
			name:     "Zero as string",
			json:     `{"value":"0"}`,
			expected: 0,
		},
		{
			name:     "Zero as number",
			json:     `{"value":0}`,
			expected: 0,
		},
		{
			name:     "Large timestamp as string",
			json:     `{"value":"1735689600"}`, // 2025-01-01 00:00:00 UTC
			expected: 1735689600,
		},
		{
			name:     "Large timestamp as number",
			json:     `{"value":1735689600}`,
			expected: 1735689600,
		},
		{
			name:     "Negative value as string",
			json:     `{"value":"-100"}`,
			expected: -100,
		},
		{
			name:     "Negative value as number",
			json:     `{"value":-100}`,
			expected: -100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result struct {
				Value zabbix.StringInt64 `json:"value"`
			}

			err := json.Unmarshal([]byte(tc.json), &result)
			require.NoError(t, err, "Failed to unmarshal")
			assert.Equal(t, tc.expected, result.Value.Int64(), "Value mismatch")
		})
	}
}

// TestBoolStringEdgeCases tests BoolString handling of edge cases
func TestBoolStringEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected bool
		shouldError bool
	}{
		{
			name:     "String zero",
			json:     `{"value":"0"}`,
			expected: false,
		},
		{
			name:     "String one",
			json:     `{"value":"1"}`,
			expected: true,
		},
		{
			name:     "Boolean true",
			json:     `{"value":true}`,
			expected: true,
		},
		{
			name:     "Boolean false",
			json:     `{"value":false}`,
			expected: false,
		},
		{
			name:        "Invalid string value",
			json:        `{"value":"2"}`,
			shouldError: true,
		},
		{
			name:        "Invalid string true",
			json:        `{"value":"true"}`,
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result struct {
				Value zabbix.BoolString `json:"value"`
			}

			err := json.Unmarshal([]byte(tc.json), &result)

			if tc.shouldError {
				assert.Error(t, err, "Expected error for invalid value")
			} else {
				require.NoError(t, err, "Failed to unmarshal")
				assert.Equal(t, tc.expected, result.Value.Bool(), "Value mismatch")
			}
		})
	}
}

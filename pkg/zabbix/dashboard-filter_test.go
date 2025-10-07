package zabbix

import (
	"testing"
)

func TestParseProblemsWidgetFilters_NilDashboard(t *testing.T) {
	_, err := ParseProblemsWidgetFilters(nil)
	if err == nil {
		t.Error("Expected error for nil dashboard, got nil")
	}
}

func TestParseProblemsWidgetFilters_NoProblemsWidget(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{Type: "graph"},
					{Type: "map"},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(options) != 0 {
		t.Errorf("Expected 0 options for dashboard without problems widget, got %d", len(options))
	}
}

func TestParseProblemsWidgetFilters_EmptyWidget(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{
						Type:   "problems",
						Fields: []WidgetField{},
					},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(options) != 0 {
		t.Errorf("Expected 0 options for empty widget, got %d", len(options))
	}
}

func TestParseProblemsWidgetFilters_GroupIDs(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{
						Type: "problems",
						Fields: []WidgetField{
							{Type: "2", Name: "groupids", Value: "4"},
							{Type: "2", Name: "groupids", Value: "5"},
						},
					},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(options) == 0 {
		t.Error("Expected options for groupids filter")
	}
}

func TestParseProblemsWidgetFilters_Severities(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{
						Type: "problems",
						Fields: []WidgetField{
							{Type: "3", Name: "severities", Value: "4"},
							{Type: "3", Name: "severities", Value: "5"},
						},
					},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(options) == 0 {
		t.Error("Expected options for severities filter")
	}
}

func TestParseProblemsWidgetFilters_Tags(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{
						Type: "problems",
						Fields: []WidgetField{
							{Type: "1", Name: "tags.tag.0", Value: "Service"},
							{Type: "1", Name: "tags.value.0", Value: "Database"},
							{Type: "0", Name: "tags.operator.0", Value: "0"},
							{Type: "1", Name: "tags.tag.1", Value: "Environment"},
							{Type: "1", Name: "tags.value.1", Value: "Production"},
							{Type: "0", Name: "tags.operator.1", Value: "1"},
						},
					},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(options) == 0 {
		t.Error("Expected options for tags filter")
	}
}

func TestParseProblemsWidgetFilters_MultipleFilters(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{
						Type: "problems",
						Fields: []WidgetField{
							{Type: "2", Name: "groupids", Value: "4"},
							{Type: "2", Name: "hostids", Value: "10084"},
							{Type: "3", Name: "severities", Value: "4"},
							{Type: "3", Name: "severities", Value: "5"},
							{Type: "0", Name: "show_suppressed", Value: "1"},
							{Type: "0", Name: "unacknowledged", Value: "1"},
						},
					},
				},
			},
		},
	}

	options, err := ParseProblemsWidgetFilters(dashboard)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Should have at least: groupids, hostids, severities, show_suppressed, unacknowledged
	if len(options) < 5 {
		t.Errorf("Expected at least 5 options, got %d", len(options))
	}
}

func TestFindProblemsWidget_FirstWidget(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{Type: "problems", Name: "First"},
					{Type: "problems", Name: "Second"},
				},
			},
		},
	}

	widget := findProblemsWidget(dashboard)
	if widget == nil {
		t.Fatal("Expected to find problems widget")
	}
	if widget.Name != "First" {
		t.Errorf("Expected to find first problems widget, got %s", widget.Name)
	}
}

func TestFindProblemsWidget_MultiplePages(t *testing.T) {
	dashboard := &Dashboard{
		Pages: []DashboardPage{
			{
				Widgets: []Widget{
					{Type: "graph"},
				},
			},
			{
				Widgets: []Widget{
					{Type: "problems", Name: "OnSecondPage"},
				},
			},
		},
	}

	widget := findProblemsWidget(dashboard)
	if widget == nil {
		t.Fatal("Expected to find problems widget")
	}
	if widget.Name != "OnSecondPage" {
		t.Errorf("Expected to find problems widget on second page, got %s", widget.Name)
	}
}

func TestExtractTags_ValidTags(t *testing.T) {
	fields := []WidgetField{
		{Type: "1", Name: "tags.tag.0", Value: "Service"},
		{Type: "1", Name: "tags.value.0", Value: "Database"},
		{Type: "0", Name: "tags.operator.0", Value: "0"},
		{Type: "1", Name: "tags.tag.1", Value: "Environment"},
		{Type: "1", Name: "tags.value.1", Value: "Production"},
		{Type: "0", Name: "tags.operator.1", Value: "1"},
	}

	tags := extractTags(fields)
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	// Check first tag
	foundTag0 := false
	for _, tag := range tags {
		if tag.Tag == "Service" && tag.Value == "Database" && tag.Operator == "0" {
			foundTag0 = true
			break
		}
	}
	if !foundTag0 {
		t.Error("Expected to find Service=Database tag with operator 0")
	}

	// Check second tag
	foundTag1 := false
	for _, tag := range tags {
		if tag.Tag == "Environment" && tag.Value == "Production" && tag.Operator == "1" {
			foundTag1 = true
			break
		}
	}
	if !foundTag1 {
		t.Error("Expected to find Environment=Production tag with operator 1")
	}
}

func TestExtractTags_IncompleteTag(t *testing.T) {
	fields := []WidgetField{
		{Type: "1", Name: "tags.tag.0", Value: "Service"},
		// Missing value and operator for index 0
		{Type: "1", Name: "tags.tag.1", Value: "Environment"},
		{Type: "1", Name: "tags.value.1", Value: "Production"},
	}

	tags := extractTags(fields)
	// Should include both tags even if incomplete
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestExtractTags_NoTags(t *testing.T) {
	fields := []WidgetField{
		{Type: "2", Name: "groupids", Value: "4"},
		{Type: "3", Name: "severities", Value: "5"},
	}

	tags := extractTags(fields)
	if len(tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(tags))
	}
}

func TestMergeProblemOptions(t *testing.T) {
	priority := []GetProblemOption{
		GetProblemOptionSeverities([]string{"5"}),
	}
	dashboard := []GetProblemOption{
		GetProblemOptionGroupsIDs([]string{"4"}),
		GetProblemOptionSeverities([]string{"3", "4"}),
	}

	merged := MergeProblemOptions(priority, dashboard)

	// Should have 3 options total (1 from priority + 2 from dashboard)
	if len(merged) != 3 {
		t.Errorf("Expected 3 merged options, got %d", len(merged))
	}
}

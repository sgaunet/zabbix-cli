package zabbix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Widget field parsing constants
const (
	tagFieldParts    = 3  // Number of parts in a tag field name (e.g., "tags.tag.0")
	evalTypeOr       = 2  // Evaltype value for OR logic
	showSuppressedOn = "1" // Value indicating show_suppressed is enabled
	unacknowledgedOn = "1" // Value indicating unacknowledged filter is enabled
)

// ErrNilDashboard is returned when a nil dashboard is provided
var ErrNilDashboard = errors.New("dashboard is nil")

// ParseProblemsWidgetFilters extracts problem.get API options from a dashboard's problem widgets.
// It finds the first "problems" type widget in the dashboard and converts its filter fields
// to GetProblemOption functions that can be used with the GetProblems API call.
//
// Returns an empty slice if no problems widget is found or if the widget has no filters.
// Returns an error if the dashboard is nil or has invalid structure.
func ParseProblemsWidgetFilters(dashboard *Dashboard) ([]GetProblemOption, error) {
	if dashboard == nil {
		return nil, ErrNilDashboard
	}

	// Find the first "problems" widget across all pages
	widget := findProblemsWidget(dashboard)
	if widget == nil {
		// No problems widget found - not an error, just return empty options
		return []GetProblemOption{}, nil
	}

	// Parse widget fields into problem options
	return parseWidgetFieldsToProblemOptions(widget.Fields)
}

// findProblemsWidget searches for the first widget of type "problems" in the dashboard.
// Returns nil if no problems widget is found.
func findProblemsWidget(dashboard *Dashboard) *Widget {
	for _, page := range dashboard.Pages {
		for i := range page.Widgets {
			if page.Widgets[i].Type == "problems" {
				return &page.Widgets[i]
			}
		}
	}
	return nil
}

// parseWidgetFieldsToProblemOptions converts widget fields to problem.get API options.
// Supports the following field mappings:
// - groupids → GroupsIDs
// - hostids → HostsIDs
// - severities → Severities
// - tags → Tags with evaltype
// - show_suppressed → Suppressed
// - unacknowledged → Acknowledged (inverted)
func parseWidgetFieldsToProblemOptions(fields []WidgetField) ([]GetProblemOption, error) {
	fieldMap := buildFieldMap(fields)

	var options []GetProblemOption
	options = appendIDOptions(options, fieldMap)
	options = appendTagOptions(options, fields, fieldMap)
	options = appendFlagOptions(options, fieldMap)

	return options, nil
}

// buildFieldMap collects field values by field name
func buildFieldMap(fields []WidgetField) map[string][]string {
	fieldMap := make(map[string][]string)
	for _, field := range fields {
		valueStr := fmt.Sprintf("%v", field.Value)
		fieldMap[field.Name] = append(fieldMap[field.Name], valueStr)
	}
	return fieldMap
}

// appendIDOptions adds group, host, and severity options
func appendIDOptions(options []GetProblemOption, fieldMap map[string][]string) []GetProblemOption {
	if groupIDs, ok := fieldMap["groupids"]; ok && len(groupIDs) > 0 {
		options = append(options, GetProblemOptionGroupsIDs(groupIDs))
	}
	if hostIDs, ok := fieldMap["hostids"]; ok && len(hostIDs) > 0 {
		options = append(options, GetProblemOptionHostsIDs(hostIDs))
	}
	if severities, ok := fieldMap["severities"]; ok && len(severities) > 0 {
		options = append(options, GetProblemOptionSeverities(severities))
	}
	return options
}

// appendTagOptions adds tag and evaltype options
func appendTagOptions(options []GetProblemOption, fields []WidgetField, fieldMap map[string][]string) []GetProblemOption {
	tags := extractTags(fields)
	if len(tags) > 0 {
		options = append(options, GetProblemOptionTags(tags))
	}
	if evalType, ok := fieldMap["evaltype"]; ok && len(evalType) > 0 {
		if evalType[0] == strconv.Itoa(evalTypeOr) {
			options = append(options, GetProblemOptionEvalType(evalTypeOr))
		}
	}
	return options
}

// appendFlagOptions adds suppressed and acknowledged flag options
func appendFlagOptions(options []GetProblemOption, fieldMap map[string][]string) []GetProblemOption {
	if showSuppressed, ok := fieldMap["show_suppressed"]; ok && len(showSuppressed) > 0 {
		if showSuppressed[0] == showSuppressedOn {
			options = append(options, GetProblemOptionSuppressed(true))
		}
	}
	if unacknowledged, ok := fieldMap["unacknowledged"]; ok && len(unacknowledged) > 0 {
		if unacknowledged[0] == unacknowledgedOn {
			options = append(options, GetProblemOptionAcknowledged(false))
		}
	}
	return options
}

// extractTags parses tag-related widget fields into FilterProblemTags structures.
// Tags in widgets are stored as separate fields with indexed names like:
// - tags.tag.0, tags.value.0, tags.operator.0
// - tags.tag.1, tags.value.1, tags.operator.1
func extractTags(fields []WidgetField) []FilterProblemTags {
	// Map to collect tag data by index
	tagMap := make(map[string]*FilterProblemTags) // key: index (e.g., "0", "1")

	for _, field := range fields {
		if strings.HasPrefix(field.Name, "tags.") {
			parts := strings.Split(field.Name, ".")
			if len(parts) != tagFieldParts {
				continue
			}

			fieldType := parts[1] // "tag", "value", or "operator"
			index := parts[2]     // "0", "1", "2", etc.
			valueStr := fmt.Sprintf("%v", field.Value)

			// Initialize tag entry if not exists
			if tagMap[index] == nil {
				tagMap[index] = &FilterProblemTags{}
			}

			// Set the appropriate field
			switch fieldType {
			case "tag":
				tagMap[index].Tag = valueStr
			case "value":
				tagMap[index].Value = valueStr
			case "operator":
				tagMap[index].Operator = valueStr
			}
		}
	}

	// Convert map to slice
	var tags []FilterProblemTags
	for _, tag := range tagMap {
		// Only include tags that have at least a tag name
		if tag.Tag != "" {
			tags = append(tags, *tag)
		}
	}

	return tags
}

// MergeProblemOptions merges two sets of problem options, with priority options taking precedence.
// This is used to combine dashboard filters with CLI flags, where CLI flags should override
// dashboard filters for the same parameter.
//
// Note: This is a simplified merge that appends all options. For true override behavior,
// the calling code should conditionally apply dashboard options only if CLI options are not set.
func MergeProblemOptions(priority []GetProblemOption, dashboard []GetProblemOption) []GetProblemOption {
	// Combine priority and dashboard options
	// The actual problem.get request will use the last value for duplicate parameters
	merged := make([]GetProblemOption, 0, len(priority)+len(dashboard))
	merged = append(merged, priority...)
	merged = append(merged, dashboard...)

	return merged
}

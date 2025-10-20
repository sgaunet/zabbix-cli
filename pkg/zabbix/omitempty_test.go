package zabbix_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

// TestConfigurationExportParamsOmitEmpty verifies that empty Format field is omitted from JSON
func TestConfigurationExportParamsOmitEmpty(t *testing.T) {
	params := zabbix.ConfigurationExportParams{
		Options: zabbix.ConfigurationExportOptions{
			TemplatesID: []string{"10001"},
		},
		// Format is intentionally left empty to test omitempty
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal ConfigurationExportParams: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "format" field is NOT present in JSON when empty
	if strings.Contains(jsonStr, `"format"`) {
		t.Errorf("expected 'format' field to be omitted when empty, but found it in JSON: %s", jsonStr)
	}

	// Verify that "options" field IS present
	if !strings.Contains(jsonStr, `"options"`) {
		t.Errorf("expected 'options' field to be present in JSON, but it's missing: %s", jsonStr)
	}
}

// TestConfigurationExportParamsIncludesFormat verifies that non-empty Format field is included
func TestConfigurationExportParamsIncludesFormat(t *testing.T) {
	params := zabbix.ConfigurationExportParams{
		Options: zabbix.ConfigurationExportOptions{
			TemplatesID: []string{"10001"},
		},
		Format: "json",
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal ConfigurationExportParams: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "format" field IS present when set
	if !strings.Contains(jsonStr, `"format":"json"`) {
		t.Errorf("expected 'format' field to be present when set, but it's missing in JSON: %s", jsonStr)
	}
}

// TestEventsAcknowledgeParamsOmitEmpty verifies that empty Message and Severity are omitted
func TestEventsAcknowledgeParamsOmitEmpty(t *testing.T) {
	params := zabbix.EventsAcknowledgeParams{
		Eventids: []string{"12345"},
		Action:   2, // acknowledge event
		// Message and Severity intentionally left empty to test omitempty
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal EventsAcknowledgeParams: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "message" field is NOT present when empty
	if strings.Contains(jsonStr, `"message"`) {
		t.Errorf("expected 'message' field to be omitted when empty, but found it in JSON: %s", jsonStr)
	}

	// Verify that "severity" field is NOT present when empty (zero value)
	if strings.Contains(jsonStr, `"severity"`) {
		t.Errorf("expected 'severity' field to be omitted when empty, but found it in JSON: %s", jsonStr)
	}

	// Verify required fields ARE present
	if !strings.Contains(jsonStr, `"eventids"`) {
		t.Errorf("expected 'eventids' field to be present, but it's missing: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"action"`) {
		t.Errorf("expected 'action' field to be present, but it's missing: %s", jsonStr)
	}
}

// TestEventsAcknowledgeParamsIncludesMessage verifies that non-empty Message is included
func TestEventsAcknowledgeParamsIncludesMessage(t *testing.T) {
	params := zabbix.EventsAcknowledgeParams{
		Eventids: []string{"12345"},
		Action:   6, // acknowledge event (2) + add message (4)
		Message:  "Issue acknowledged and being investigated",
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal EventsAcknowledgeParams: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "message" field IS present when set
	if !strings.Contains(jsonStr, `"message":"Issue acknowledged and being investigated"`) {
		t.Errorf("expected 'message' field to be present when set, but it's missing in JSON: %s", jsonStr)
	}
}

// TestEventsAcknowledgeParamsIncludesSeverity verifies that non-zero Severity is included
func TestEventsAcknowledgeParamsIncludesSeverity(t *testing.T) {
	params := zabbix.EventsAcknowledgeParams{
		Eventids: []string{"12345"},
		Action:   10, // acknowledge event (2) + change severity (8)
		Severity: 3,  // Average severity
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal EventsAcknowledgeParams: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "severity" field IS present when set
	if !strings.Contains(jsonStr, `"severity":3`) {
		t.Errorf("expected 'severity' field to be present when set, but it's missing in JSON: %s", jsonStr)
	}
}

// TestParamsImportOmitEmpty verifies that empty Format and Rules fields are omitted
// Note: This test uses a mock struct since paramsImport is not exported
func TestParamsImportOmitEmpty(t *testing.T) {
	// Create a test struct that mirrors paramsImport structure
	type testParamsImport struct {
		Format string `json:"format,omitempty"`
		Rules  string `json:"rules,omitempty"`
		Source string `json:"source"`
	}

	params := testParamsImport{
		Source: "zabbix_export:\n  version: '7.2'",
		// Format and Rules intentionally left empty to test omitempty
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("failed to marshal testParamsImport: %v", err)
	}

	jsonStr := string(jsonData)

	// Verify that "format" field is NOT present when empty
	if strings.Contains(jsonStr, `"format"`) {
		t.Errorf("expected 'format' field to be omitted when empty, but found it in JSON: %s", jsonStr)
	}

	// Verify that "rules" field is NOT present when empty
	if strings.Contains(jsonStr, `"rules"`) {
		t.Errorf("expected 'rules' field to be omitted when empty, but found it in JSON: %s", jsonStr)
	}

	// Verify that "source" field IS present
	if !strings.Contains(jsonStr, `"source"`) {
		t.Errorf("expected 'source' field to be present, but it's missing: %s", jsonStr)
	}
}

// TestOmitEmptyRoundTrip verifies that omitted fields unmarshal correctly as zero values
func TestOmitEmptyRoundTrip(t *testing.T) {
	// Test with ConfigurationExportParams
	original := zabbix.ConfigurationExportParams{
		Options: zabbix.ConfigurationExportOptions{
			TemplatesID: []string{"10001"},
		},
		// Format intentionally omitted
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var result zabbix.ConfigurationExportParams
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify Format is zero value (empty string) after round trip
	if result.Format != "" {
		t.Errorf("expected Format to be empty string after round trip, got: %s", result.Format)
	}

	// Verify Options are preserved
	if len(result.Options.TemplatesID) != 1 || result.Options.TemplatesID[0] != "10001" {
		t.Errorf("expected TemplatesID to be preserved, got: %v", result.Options.TemplatesID)
	}
}

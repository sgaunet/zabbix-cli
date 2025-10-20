package zabbix_test

import (
	"encoding/json"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

// TestBoolStringUnmarshalFromStringTrue tests unmarshaling from "1" string
func TestBoolStringUnmarshalFromStringTrue(t *testing.T) {
	jsonData := `{"value":"1"}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if !result.Value.Bool() {
		t.Errorf("expected true, got false")
	}
}

// TestBoolStringUnmarshalFromStringFalse tests unmarshaling from "0" string
func TestBoolStringUnmarshalFromStringFalse(t *testing.T) {
	jsonData := `{"value":"0"}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if result.Value.Bool() {
		t.Errorf("expected false, got true")
	}
}

// TestBoolStringUnmarshalFromBoolTrue tests unmarshaling from boolean true
func TestBoolStringUnmarshalFromBoolTrue(t *testing.T) {
	jsonData := `{"value":true}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if !result.Value.Bool() {
		t.Errorf("expected true, got false")
	}
}

// TestBoolStringUnmarshalFromBoolFalse tests unmarshaling from boolean false
func TestBoolStringUnmarshalFromBoolFalse(t *testing.T) {
	jsonData := `{"value":false}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if result.Value.Bool() {
		t.Errorf("expected false, got true")
	}
}

// TestBoolStringUnmarshalInvalidString tests unmarshaling from invalid string
func TestBoolStringUnmarshalInvalidString(t *testing.T) {
	jsonData := `{"value":"2"}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err == nil {
		t.Errorf("expected error for invalid value '2', got nil")
	}
}

// TestBoolStringUnmarshalInvalidType tests unmarshaling from invalid type
func TestBoolStringUnmarshalInvalidType(t *testing.T) {
	jsonData := `{"value":42}`
	var result struct {
		Value zabbix.BoolString `json:"value"`
	}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err == nil {
		t.Errorf("expected error for invalid type (number), got nil")
	}
}

// TestBoolStringMarshalTrue tests marshaling true to "1"
func TestBoolStringMarshalTrue(t *testing.T) {
	data := struct {
		Value zabbix.BoolString `json:"value"`
	}{
		Value: zabbix.BoolString(true),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}

	expected := `{"value":"1"}`
	if string(jsonData) != expected {
		t.Errorf("expected %s, got %s", expected, string(jsonData))
	}
}

// TestBoolStringMarshalFalse tests marshaling false to "0"
func TestBoolStringMarshalFalse(t *testing.T) {
	data := struct {
		Value zabbix.BoolString `json:"value"`
	}{
		Value: zabbix.BoolString(false),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}

	expected := `{"value":"0"}`
	if string(jsonData) != expected {
		t.Errorf("expected %s, got %s", expected, string(jsonData))
	}
}

// TestBoolStringBoolMethod tests the Bool() conversion method
func TestBoolStringBoolMethod(t *testing.T) {
	bsTrue := zabbix.BoolString(true)
	if !bsTrue.Bool() {
		t.Errorf("expected true, got false")
	}

	bsFalse := zabbix.BoolString(false)
	if bsFalse.Bool() {
		t.Errorf("expected false, got true")
	}
}

// TestBoolStringRoundTrip tests complete marshal/unmarshal cycle
func TestBoolStringRoundTrip(t *testing.T) {
	original := struct {
		Enabled  zabbix.BoolString `json:"enabled"`
		Disabled zabbix.BoolString `json:"disabled"`
	}{
		Enabled:  zabbix.BoolString(true),
		Disabled: zabbix.BoolString(false),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var result struct {
		Enabled  zabbix.BoolString `json:"enabled"`
		Disabled zabbix.BoolString `json:"disabled"`
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	// Verify values match
	if result.Enabled.Bool() != original.Enabled.Bool() {
		t.Errorf("enabled mismatch: expected %v, got %v", original.Enabled.Bool(), result.Enabled.Bool())
	}
	if result.Disabled.Bool() != original.Disabled.Bool() {
		t.Errorf("disabled mismatch: expected %v, got %v", original.Disabled.Bool(), result.Disabled.Bool())
	}
}

// TestBoolStringInDashboardStruct tests BoolString usage in Dashboard struct
func TestBoolStringInDashboardStruct(t *testing.T) {
	jsonData := `{
		"dashboardid": "1",
		"name": "Test Dashboard",
		"private": "1",
		"auto_start": "0"
	}`

	var dashboard zabbix.Dashboard
	err := json.Unmarshal([]byte(jsonData), &dashboard)
	if err != nil {
		t.Errorf("failed to unmarshal dashboard: %v", err)
	}

	if !dashboard.Private.Bool() {
		t.Errorf("expected private=true, got false")
	}
	if dashboard.AutoStart.Bool() {
		t.Errorf("expected auto_start=false, got true")
	}
}

// TestBoolStringInProblemStruct tests BoolString usage in Problem struct
func TestBoolStringInProblemStruct(t *testing.T) {
	jsonData := `{
		"eventid": "123",
		"source": "0",
		"object": "0",
		"objectid": "456",
		"clock": "1234567890",
		"ns": "0",
		"name": "Test Problem",
		"acknowledged": "1",
		"severity": "3",
		"suppressed": "0"
	}`

	var problem zabbix.Problem
	err := json.Unmarshal([]byte(jsonData), &problem)
	if err != nil {
		t.Errorf("failed to unmarshal problem: %v", err)
	}

	if !problem.Acknowledged.Bool() {
		t.Errorf("expected acknowledged=true, got false")
	}
	if problem.Suppressed.Bool() {
		t.Errorf("expected suppressed=false, got true")
	}

	// Test helper methods
	if !problem.GetAcknowledge() {
		t.Errorf("GetAcknowledge() expected true, got false")
	}
	if problem.GetSuppressed() {
		t.Errorf("GetSuppressed() expected false, got true")
	}
}

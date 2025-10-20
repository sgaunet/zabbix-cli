package zabbix

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FormatType represents the configuration data format
type FormatType string

const (
	// FormatYAML represents YAML format
	FormatYAML FormatType = "yaml"
	// FormatJSON represents JSON format
	FormatJSON FormatType = "json"
	// FormatXML represents XML format
	FormatXML FormatType = "xml"
	// FormatUnknown represents unknown format
	FormatUnknown FormatType = "unknown"
)

// DetectFormatFromExtension detects format from file extension
func DetectFormatFromExtension(filename string) FormatType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	case ".json":
		return FormatJSON
	case ".xml":
		return FormatXML
	default:
		return FormatUnknown
	}
}

// DetectFormatFromContent attempts to detect format from content
func DetectFormatFromContent(data string) FormatType {
	trimmed := strings.TrimSpace(data)
	if trimmed == "" {
		return FormatUnknown
	}

	// Check for JSON - starts with { or [
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(data), &js); err == nil {
			return FormatJSON
		}
	}

	// Check for XML - starts with <
	if strings.HasPrefix(trimmed, "<") || strings.HasPrefix(trimmed, "<?xml") {
		var xmlData any
		if err := xml.Unmarshal([]byte(data), &xmlData); err == nil {
			return FormatXML
		}
	}

	// Try YAML - most permissive format
	var yamlData any
	if err := yaml.Unmarshal([]byte(data), &yamlData); err == nil {
		return FormatYAML
	}

	return FormatUnknown
}

// ValidateFormat validates that the given data is valid for the specified format
func ValidateFormat(data string, format FormatType) error {
	if data == "" {
		return fmt.Errorf("empty data")
	}

	switch format {
	case FormatYAML:
		return ValidateYAML(data)
	case FormatJSON:
		return ValidateJSON(data)
	case FormatXML:
		return ValidateXML(data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// ValidateYAML validates YAML format
func ValidateYAML(data string) error {
	var yamlData any
	if err := yaml.Unmarshal([]byte(data), &yamlData); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}
	return nil
}

// ValidateJSON validates JSON format
func ValidateJSON(data string) error {
	var jsonData any
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	return nil
}

// ValidateXML validates XML format
func ValidateXML(data string) error {
	var xmlData any
	if err := xml.Unmarshal([]byte(data), &xmlData); err != nil {
		return fmt.Errorf("invalid XML format: %w", err)
	}
	return nil
}

// ValidateZabbixExportData validates that data looks like valid Zabbix export data
func ValidateZabbixExportData(data string) error {
	if data == "" {
		return fmt.Errorf("empty export data")
	}

	// Check for zabbix_export marker (common to all formats)
	if !strings.Contains(data, "zabbix_export") {
		return fmt.Errorf("data does not appear to be Zabbix export data (missing 'zabbix_export' marker)")
	}

	return nil
}

// String returns the string representation of FormatType
func (f FormatType) String() string {
	return string(f)
}

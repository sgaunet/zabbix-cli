package zabbix_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestDetectFormatFromExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     zabbix.FormatType
	}{
		{"YAML with .yaml extension", "template.yaml", zabbix.FormatYAML},
		{"YAML with .yml extension", "template.yml", zabbix.FormatYAML},
		{"JSON with .json extension", "template.json", zabbix.FormatJSON},
		{"XML with .xml extension", "template.xml", zabbix.FormatXML},
		{"Unknown extension .txt", "template.txt", zabbix.FormatUnknown},
		{"No extension", "template", zabbix.FormatUnknown},
		{"Mixed case .YAML", "template.YAML", zabbix.FormatYAML},
		{"Mixed case .JSON", "template.JSON", zabbix.FormatJSON},
		{"Full path YAML", "/path/to/template.yaml", zabbix.FormatYAML},
		{"Full path JSON", "/path/to/template.json", zabbix.FormatJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := zabbix.DetectFormatFromExtension(tt.filename)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDetectFormatFromContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    zabbix.FormatType
	}{
		{
			name: "Valid JSON object",
			content: `{
				"zabbix_export": {
					"version": "7.2"
				}
			}`,
			want: zabbix.FormatJSON,
		},
		{
			name:    "Valid JSON array",
			content: `["item1", "item2"]`,
			want:    zabbix.FormatJSON,
		},
		{
			name: "Valid XML",
			content: `<?xml version="1.0" encoding="UTF-8"?>
			<zabbix_export>
				<version>7.2</version>
			</zabbix_export>`,
			want: zabbix.FormatXML,
		},
		{
			name: "Valid XML without declaration",
			content: `<zabbix_export>
				<version>7.2</version>
			</zabbix_export>`,
			want: zabbix.FormatXML,
		},
		{
			name: "Valid YAML",
			content: `zabbix_export:
  version: '7.2'
  templates:
    - name: Test`,
			want: zabbix.FormatYAML,
		},
		{
			name:    "Empty string",
			content: "",
			want:    zabbix.FormatUnknown,
		},
		{
			name:    "Whitespace only",
			content: "   \n  \t  ",
			want:    zabbix.FormatUnknown,
		},
		{
			name:    "Invalid format",
			content: "not valid format at all",
			want:    zabbix.FormatYAML, // YAML is most permissive and accepts plain strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := zabbix.DetectFormatFromContent(tt.content)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "Valid YAML",
			data: `zabbix_export:
  version: '7.2'
  templates:
    - name: Test`,
			wantErr: false,
		},
		{
			name:    "Simple key-value YAML",
			data:    "key: value",
			wantErr: false,
		},
		{
			name: "Invalid YAML - bad indentation",
			data: `zabbix_export:
version: '7.2'
  templates:`,
			wantErr: true,
		},
		{
			name:    "Invalid YAML - unclosed quote",
			data:    `key: "unclosed`,
			wantErr: true,
		},
		{
			name: "Invalid YAML - tabs",
			data: "key:\n\tvalue",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := zabbix.ValidateYAML(tt.data)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid YAML format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid JSON object",
			data:    `{"key": "value"}`,
			wantErr: false,
		},
		{
			name: "Valid JSON nested object",
			data: `{
				"zabbix_export": {
					"version": "7.2",
					"templates": []
				}
			}`,
			wantErr: false,
		},
		{
			name:    "Valid JSON array",
			data:    `["item1", "item2"]`,
			wantErr: false,
		},
		{
			name:    "Invalid JSON - missing closing brace",
			data:    `{"key": "value"`,
			wantErr: true,
		},
		{
			name:    "Invalid JSON - trailing comma",
			data:    `{"key": "value",}`,
			wantErr: true,
		},
		{
			name:    "Invalid JSON - single quotes",
			data:    `{'key': 'value'}`,
			wantErr: true,
		},
		{
			name:    "Invalid JSON - not JSON",
			data:    `key: value`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := zabbix.ValidateJSON(tt.data)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid JSON format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "Valid XML with declaration",
			data: `<?xml version="1.0" encoding="UTF-8"?>
			<zabbix_export>
				<version>7.2</version>
			</zabbix_export>`,
			wantErr: false,
		},
		{
			name: "Valid XML without declaration",
			data: `<zabbix_export>
				<version>7.2</version>
			</zabbix_export>`,
			wantErr: false,
		},
		{
			name:    "Valid XML single tag",
			data:    `<root></root>`,
			wantErr: false,
		},
		{
			name:    "Invalid XML - unclosed tag",
			data:    `<root><child></root>`,
			wantErr: true,
		},
		{
			name:    "Invalid XML - no closing tag",
			data:    `<root><child>`,
			wantErr: true,
		},
		{
			name:    "Invalid XML - malformed",
			data:    `<root>><child</root>`,
			wantErr: true,
		},
		{
			name:    "Invalid XML - not XML",
			data:    `key: value`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := zabbix.ValidateXML(tt.data)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid XML format")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		format  zabbix.FormatType
		wantErr bool
	}{
		{
			name:    "Valid YAML",
			data:    "key: value",
			format:  zabbix.FormatYAML,
			wantErr: false,
		},
		{
			name:    "Invalid YAML",
			data:    "key:\n\tvalue",
			format:  zabbix.FormatYAML,
			wantErr: true,
		},
		{
			name:    "Valid JSON",
			data:    `{"key": "value"}`,
			format:  zabbix.FormatJSON,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			data:    `{key: value}`,
			format:  zabbix.FormatJSON,
			wantErr: true,
		},
		{
			name:    "Valid XML",
			data:    `<root></root>`,
			format:  zabbix.FormatXML,
			wantErr: false,
		},
		{
			name:    "Invalid XML",
			data:    `<root>`,
			format:  zabbix.FormatXML,
			wantErr: true,
		},
		{
			name:    "Empty data",
			data:    "",
			format:  zabbix.FormatYAML,
			wantErr: true,
		},
		{
			name:    "Unknown format",
			data:    "some data",
			format:  zabbix.FormatUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := zabbix.ValidateFormat(tt.data, tt.format)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateZabbixExportData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "Valid Zabbix YAML export",
			data: `zabbix_export:
  version: '7.2'
  templates: []`,
			wantErr: false,
		},
		{
			name: "Valid Zabbix JSON export",
			data: `{
				"zabbix_export": {
					"version": "7.2"
				}
			}`,
			wantErr: false,
		},
		{
			name: "Valid Zabbix XML export",
			data: `<?xml version="1.0" encoding="UTF-8"?>
			<zabbix_export>
				<version>7.2</version>
			</zabbix_export>`,
			wantErr: false,
		},
		{
			name:    "Invalid - missing zabbix_export marker",
			data:    `templates: []`,
			wantErr: true,
		},
		{
			name:    "Invalid - empty data",
			data:    "",
			wantErr: true,
		},
		{
			name:    "Invalid - wrong format",
			data:    `random: data`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := zabbix.ValidateZabbixExportData(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFormatTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format zabbix.FormatType
		want   string
	}{
		{"YAML format", zabbix.FormatYAML, "yaml"},
		{"JSON format", zabbix.FormatJSON, "json"},
		{"XML format", zabbix.FormatXML, "xml"},
		{"Unknown format", zabbix.FormatUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.format.String()
			require.Equal(t, tt.want, got)
		})
	}
}

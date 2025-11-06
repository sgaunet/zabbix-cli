package cmd_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/cmd"
)

func TestGetFormatOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		format      string
		shouldError bool
	}{
		{
			name:        "Valid YAML format lowercase",
			format:      "yaml",
			shouldError: false,
		},
		{
			name:        "Valid YAML format uppercase",
			format:      "YAML",
			shouldError: false,
		},
		{
			name:        "Valid JSON format lowercase",
			format:      "json",
			shouldError: false,
		},
		{
			name:        "Valid JSON format uppercase",
			format:      "JSON",
			shouldError: false,
		},
		{
			name:        "Valid XML format lowercase",
			format:      "xml",
			shouldError: false,
		},
		{
			name:        "Valid XML format uppercase",
			format:      "XML",
			shouldError: false,
		},
		{
			name:        "Invalid format",
			format:      "invalid",
			shouldError: true,
		},
		{
			name:        "Empty format",
			format:      "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opt, err := cmd.GetFormatOption(tt.format)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error for format %s, got nil", tt.format)
				}
				if opt != nil {
					t.Errorf("expected nil option for invalid format %s, got non-nil", tt.format)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for format %s: %v", tt.format, err)
				}
				if opt == nil {
					t.Errorf("expected non-nil option for format %s, got nil", tt.format)
				}
			}
		})
	}
}

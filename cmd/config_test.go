package cmd_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/cmd"
)

func TestPrintConfigCmd(t *testing.T) {
	t.Parallel()
	// Run the test
	cmd := cmd.PrintConfigCmd
	if cmd.Use != "config" {
		t.Errorf("expected %s, got %s", "config", cmd.Use)
	}
	if cmd.Run == nil {
		t.Errorf("expected non-nil Run, got nil")
	}
}

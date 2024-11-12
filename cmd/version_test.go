package cmd_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/cmd"
)

func TestVersionCmd(t *testing.T) {
	t.Parallel()
	// Run the test
	cmd := cmd.VersionCmd
	if cmd.Use != "version" {
		t.Errorf("expected %s, got %s", "config", cmd.Use)
	}
	if cmd.Run == nil {
		t.Errorf("expected non-nil Run, got nil")
	}
}

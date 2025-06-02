package zabbix_test

import (
	"testing"
	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		sev      zabbix.Severity
		expected string
	}{
		{zabbix.NotClassified, "Not classified"},
		{zabbix.Information, "Information"},
		{zabbix.Warning, "Warning"},
		{zabbix.Average, "Average"},
		{zabbix.High, "High"},
		{zabbix.Disaster, "Disaster"},
		{zabbix.Severity(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.sev.String(); got != c.expected {
			t.Errorf("Severity(%d).String() = %q, want %q", c.sev, got, c.expected)
		}
	}
}

func TestGetSeverity(t *testing.T) {
	cases := []struct {
		input    int
		expected zabbix.Severity
	}{
		{0, zabbix.NotClassified},
		{1, zabbix.Information},
		{2, zabbix.Warning},
		{3, zabbix.Average},
		{4, zabbix.High},
		{5, zabbix.Disaster},
		{99, zabbix.NotClassified},
	}
	for _, c := range cases {
		if got := zabbix.GetSeverity(c.input); got != c.expected {
			t.Errorf("GetSeverity(%d) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestGetSeverityString(t *testing.T) {
	cases := []struct {
		input    string
		expected zabbix.Severity
	}{
		{"Not classified", zabbix.NotClassified},
		{"Information", zabbix.Information},
		{"Warning", zabbix.Warning},
		{"Average", zabbix.Average},
		{"High", zabbix.High},
		{"Disaster", zabbix.Disaster},
		{"foobar", zabbix.NotClassified},
	}
	for _, c := range cases {
		if got := zabbix.GetSeverityString(c.input); got != c.expected {
			t.Errorf("GetSeverityString(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestNewSeverity(t *testing.T) {
	if got := zabbix.NewSeverity(3); got != zabbix.Average {
		t.Errorf("NewSeverity(3) = %v, want %v", got, zabbix.Average)
	}
}

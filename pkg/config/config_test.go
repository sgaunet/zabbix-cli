package config_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/config"
)

func TestValid(t *testing.T) {
	t.Parallel()
	c := config.Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "admin",
		ZabbixPassword: "*****",
	}
	if !c.IsValid() {
		t.Errorf("Config is valid")
	}
}

func TestInvalidMissingZabbixEndpoint(t *testing.T) {
	t.Parallel()
	c := config.Config{
		ZabbixEndpoint: "",
		ZabbixUser:     "admin",
		ZabbixPassword: "*****",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

func TestInvalidMissingZabbixUser(t *testing.T) {
	t.Parallel()
	c := config.Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "",
		ZabbixPassword: "*****",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

func TestInvalidMissingZabbixPassword(t *testing.T) {
	t.Parallel()
	c := config.Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "admin",
		ZabbixPassword: "",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

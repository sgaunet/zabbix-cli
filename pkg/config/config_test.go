package config

import "testing"

func TestValid(t *testing.T) {
	c := Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "admin",
		ZabbixPassword: "*****",
	}
	if !c.IsValid() {
		t.Errorf("Config is valid")
	}
}

func TestInvalidMissingZabbixEndpoint(t *testing.T) {
	c := Config{
		ZabbixEndpoint: "",
		ZabbixUser:     "admin",
		ZabbixPassword: "*****",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

func TestInvalidMissingZabbixUser(t *testing.T) {
	c := Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "",
		ZabbixPassword: "*****",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

func TestInvalidMissingZabbixPassword(t *testing.T) {
	c := Config{
		ZabbixEndpoint: "http://zabbix.mydomain.com/api_JSONRPC.php",
		ZabbixUser:     "admin",
		ZabbixPassword: "",
	}
	if c.IsValid() {
		t.Errorf("Config is invalid")
	}
}

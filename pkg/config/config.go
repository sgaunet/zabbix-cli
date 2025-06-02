// Package config provides configuration management for the Zabbix CLI application.
package config

// Config struct holds the configuration for the application.
type Config struct {
	ZabbixEndpoint string `mapstructure:"zabbix_endpoint"`
	ZabbixUser     string `mapstructure:"zabbix_user"`
	ZabbixPassword string `mapstructure:"zabbix_password"`
}

// IsValid checks if the configuration is valid.
func (c *Config) IsValid() bool {
	if c.ZabbixEndpoint == "" || c.ZabbixUser == "" || c.ZabbixPassword == "" {
		return false
	}
	return true
}

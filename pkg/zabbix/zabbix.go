package zabbix

import "net/http"

const JSONRPC = "2.0"

// ZabbixAPI struct holds the configuration for the Zabbix API
type ZabbixAPI struct {
	client      *http.Client
	auth        string // auth token
	APIEndpoint string
	User        string
	Password    string
	id          int // id for the JSON-RPC request - unique identifier
}

// ZbxParams struct is a part of the zbxRequestLogin struct
type ZbxParams struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// ZbxLogin struct is used to login to the Zabbix API
type ZbxRequestLogin struct {
	JSONRPC string    `json:"jsonrpc"`
	Method  string    `json:"method"`
	Params  ZbxParams `json:"params"`
	ID      int       `json:"id"`
}

// zbxRequestLogout struct is used to logout from the Zabbix API
type ZbxRequestLogout struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	Auth    string            `json:"auth"`
	ID      int               `json:"id"`
}

// ZbxLoginResponse struct is the response from the Zabbix API after a login request
type ZbxLoginResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

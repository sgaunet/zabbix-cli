package zabbix

import (
	"net/http"
	// "math/rand" // For a more robust unique ID generator later
	// "time"      // For a more robust unique ID generator later
)

// JSONRPC is the JSON-RPC version used for Zabbix API requests.
const JSONRPC = "2.0"

// Client struct holds the configuration for the Zabbix API
// Client represents the configuration for the Zabbix API.
// Client is the main struct for interacting with the Zabbix API. (Linter: stutter is intentional for public API clarity)
type Client struct {
	client      *http.Client
	auth        string // auth token
	APIEndpoint string
	User        string
	Password    string
	id          int // id for the JSON-RPC request - unique identifier
}

// Params struct is a part of the LoginRequest struct
type Params struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is used to login to the Zabbix API.
type LoginRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

// LogoutRequest is used to logout from the Zabbix API.
type LogoutRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []interface{}     `json:"params"` // Must be an empty array for logout
	Auth    string            `json:"auth"`
	ID      int               `json:"id"`
}

// LoginResponse struct is the response from the Zabbix API after a login request
type LoginResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

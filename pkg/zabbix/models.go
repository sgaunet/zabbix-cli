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
}

// zbxParams struct is a part of the zbxRequestLogin struct
type zbxParams struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// zbxLogin struct is used to login to the Zabbix API
type zbxRequestLogin struct {
	JSONRPC string    `json:"JSONRPC"`
	Method  string    `json:"method"`
	Params  zbxParams `json:"params"`
	ID      int       `json:"ID"`
}

// zbxRequestLogout struct is used to logout from the Zabbix API
type zbxRequestLogout struct {
	JSONRPC string            `json:"JSONRPC"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	Auth    string            `json:"auth"`
	ID      int               `json:"ID"`
}

type zbxLoginResponse struct {
	JSONRPC string `json:"JSONRPC"`
	Result  string `json:"result"`
	ID      int    `json:"ID"`
}

// type zbxTagsFilterProblem struct {
// 	Tag      string `json:"tag" yaml:"tag"`
// 	Value    string `json:"value" yaml:"value"`
// 	Operator string `json:"operator" yaml:"operator"`
// }

// type zbxParamsProblem struct {
// 	Suppressed   bool                   `json:"suppressed"`
// 	Recent       bool                   `json:"recent"`
// 	Acknowledged bool                   `json:"acknowledged"`
// 	TimeFrom     string                 `json:"time_from"`
// 	Tags         []zbxTagsFilterProblem `json:"tags"`
// }

// type zbxGetProblem struct {
// 	JSONRPC string           `json:"JSONRPC"`
// 	Method  string           `json:"method"`
// 	Params  zbxParamsProblem `json:"params"`
// 	Auth    string           `json:"auth"`
// 	ID      int              `json:"ID"`
// }

// type zbxProblem struct {
// 	Acknowledged  string `json:"acknowledged"`
// 	Clock         string `json:"clock"`
// 	CorrelationID string `json:"correlationID"`
// 	EventID       string `json:"eventID"`
// 	Name          string `json:"name"`
// 	Ns            string `json:"ns"`
// 	Object        string `json:"object"`
// 	ObjectID      string `json:"objectID"`
// 	Opdata        string `json:"opdata"`
// 	Rclock        string `json:"r_clock"`
// 	ReventID      string `json:"r_eventID"`
// 	Rns           string `json:"r_ns"`
// 	Severity      string `json:"severity"`
// 	Source        string `json:"source"`
// 	Suppressed    string `json:"suppressed"`
// }

// type zbxResultProblem struct {
// 	JSONRPC string       `json:"JSONRPC"`
// 	Result  []zbxProblem `json:"result"`
// 	ID      int          `json:"ID"`
// }

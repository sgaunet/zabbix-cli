package zabbix

// MaintenanceType represents the type of maintenance.
type MaintenanceType int

// Maintenance types
const (
	// MaintenanceWithDataCollection - maintenance with data collection
	MaintenanceWithDataCollection MaintenanceType = 0
	// MaintenanceNoDataCollection - maintenance without data collection
	MaintenanceNoDataCollection MaintenanceType = 1
)

// TimePeriodType represents the type of time period.
type TimePeriodType int

// Time period types
const (
	// TimePeriodTypeOneTime - one time only
	TimePeriodTypeOneTime TimePeriodType = 0
	// TimePeriodTypeDaily - daily
	TimePeriodTypeDaily TimePeriodType = 2
	// TimePeriodTypeWeekly - weekly
	TimePeriodTypeWeekly TimePeriodType = 3
	// TimePeriodTypeMonthly - monthly
	TimePeriodTypeMonthly TimePeriodType = 4
)

// Maintenance represents a Zabbix maintenance.
type Maintenance struct {
	// Maintenance ID (readonly)
	MaintenanceID string `json:"maintenanceid,omitempty"`
	// Maintenance name
	Name string `json:"name"`
	// Maintenance active since (timestamp)
	ActiveSince int64 `json:"active_since"`
	// Maintenance active till (timestamp)
	ActiveTill int64 `json:"active_till"`
	// Description of the maintenance
	Description string `json:"description"`
	// Type of maintenance
	MaintenanceType MaintenanceType `json:"maintenance_type"`
	// Time when the maintenance was created (readonly)
	CreatedAt int64 `json:"created_at,omitempty"`
	// Time when the maintenance was last updated (readonly)
	UpdatedAt int64 `json:"updated_at,omitempty"`
	// Time periods when the maintenance is active
	TimePeriods []TimePeriod `json:"timeperiods,omitempty"`
	// Host groups to add to the maintenance
	GroupIDs []string `json:"groupids,omitempty"`
	// Hosts to add to the maintenance
	HostIDs []string `json:"hostids,omitempty"`
	// Tags to filter problems during maintenance
	Tags []ProblemTag `json:"tags,omitempty"`
}

// TimePeriod represents a time period for maintenance.
type TimePeriod struct {
	// Time period ID (readonly)
	TimePeriodID string `json:"timeperiodid,omitempty"`
	// Time period type
	TimePeriodType TimePeriodType `json:"timeperiod_type"`
	// For daily and weekly periods, every day/week the maintenance will be performed starting at start_time
	Every int `json:"every,omitempty"`
	// Day of the month when the maintenance must come into effect (1-31)
	Day int `json:"day,omitempty"`
	// Days of the week when the maintenance must come into effect (1-7, where 1 is Monday)
	DayOfWeek []int `json:"dayofweek,omitempty"`
	// Start time of the maintenance period in seconds since the start of the day
	StartTime int `json:"start_time"`
	// Duration of the maintenance period in seconds
	Period int `json:"period"`
}

// ProblemTag represents a problem tag for maintenance.
type ProblemTag struct {
	// Tag name
	Tag string `json:"tag"`
	// Tag value
	Value string `json:"value,omitempty"`
	// Tag operator
	Operator int `json:"operator,omitempty"`
}

// MaintenanceResponse represents the response from maintenance API calls.
type MaintenanceResponse struct {
	MaintenanceIDs []string `json:"maintenanceids"`
}

// MaintenanceGetParams represents parameters for maintenance.get API call.
type MaintenanceGetParams struct {
	CommonGetParams
	GroupIDs          []string `json:"groupids,omitempty"`
	HostIDs           []string `json:"hostids,omitempty"`
	MaintenanceIDs    []string `json:"maintenanceids,omitempty"`
	SelectGroups      string   `json:"selectGroups,omitempty"`
	SelectHosts       string   `json:"selectHosts,omitempty"`
	SelectTimePeriods string   `json:"selectTimeperiods,omitempty"`
}

// MaintenanceCreateRequest represents a request to create a maintenance.
type MaintenanceCreateRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  Maintenance `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

// MaintenanceGetRequest represents a request to get maintenances.
type MaintenanceGetRequest struct {
	JSONRPC string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  MaintenanceGetParams `json:"params"`
	Auth    string               `json:"auth,omitempty"`
	ID      int                  `json:"id"`
}

// MaintenanceDeleteRequest represents a request to delete maintenances.
type MaintenanceDeleteRequest struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Auth    string   `json:"auth,omitempty"`
	ID      int      `json:"id"`
}

// MaintenanceUpdateRequest represents a request to update a maintenance.
type MaintenanceUpdateRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  Maintenance `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

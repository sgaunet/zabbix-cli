package zabbix

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Maintenance API method names
const (
	MethodMaintenanceCreate = "maintenance.create"
	MethodMaintenanceGet    = "maintenance.get"
	MethodMaintenanceDelete = "maintenance.delete"
)

// StringInt64 is a custom type that can unmarshal both string and integer JSON values into an int64
type StringInt64 int64

// UnmarshalJSON is a custom unmarshaler for StringInt64 to handle both string and int values
func (si *StringInt64) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as integer first
	var intValue int64
	if err := json.Unmarshal(data, &intValue); err == nil {
		*si = StringInt64(intValue)
		return nil
	}

	// If that fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return fmt.Errorf("failed to unmarshal StringInt64 from string: %w", err)
	}

	// Convert string to int64
	intValue, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse StringInt64 from string '%s': %w", stringValue, err)
	}

	*si = StringInt64(intValue)
	return nil
}

// Int64 converts StringInt64 back to a standard int64 value
func (si *StringInt64) Int64() int64 {
	if si == nil {
		return 0
	}
	return int64(*si)
}

// MaintenanceType represents the type of maintenance.
type MaintenanceType int

// Maintenance types
const (
	// MaintenanceWithDataCollection - maintenance with data collection
	MaintenanceWithDataCollection MaintenanceType = 0
	// MaintenanceNoDataCollection - maintenance without data collection
	MaintenanceNoDataCollection MaintenanceType = 1
)

// UnmarshalJSON is a custom unmarshaler for MaintenanceType to handle both string and int values
func (mt *MaintenanceType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as integer first
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*mt = MaintenanceType(intValue)
		return nil
	}

	// If that fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return fmt.Errorf("failed to unmarshal MaintenanceType from string: %w", err)
	}

	// Convert string to int
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return fmt.Errorf("failed to parse MaintenanceType from string '%s': %w", stringValue, err)
	}

	*mt = MaintenanceType(intValue)
	return nil
}

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

// UnmarshalJSON is a custom unmarshaler for TimePeriodType to handle both string and int values
func (tpt *TimePeriodType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as integer first
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*tpt = TimePeriodType(intValue)
		return nil
	}

	// If that fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return fmt.Errorf("failed to unmarshal TimePeriodType from string: %w", err)
	}

	// Convert string to int
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return fmt.Errorf("failed to parse TimePeriodType from string '%s': %w", stringValue, err)
	}

	*tpt = TimePeriodType(intValue)
	return nil
}

// Maintenance represents a Zabbix maintenance.
type Maintenance struct {
	// Maintenance ID (readonly)
	MaintenanceID string `json:"maintenanceid,omitempty"`
	// Maintenance name
	Name string `json:"name"`
	// Maintenance active since (timestamp)
	ActiveSince StringInt64 `json:"active_since"`
	// Maintenance active till (timestamp)
	ActiveTill StringInt64 `json:"active_till"`
	// Description of the maintenance
	Description string `json:"description"`
	// Type of maintenance
	MaintenanceType MaintenanceType `json:"maintenance_type"`
	// Time when the maintenance was created (readonly)
	CreatedAt StringInt64 `json:"created_at,omitempty"`
	// Time when the maintenance was last updated (readonly)
	UpdatedAt StringInt64 `json:"updated_at,omitempty"`
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
	StartTime int `json:"start_time,omitempty"`
	// Duration of the maintenance period in seconds
	Period int `json:"period,omitempty"`
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

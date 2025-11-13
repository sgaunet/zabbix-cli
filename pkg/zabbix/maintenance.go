package zabbix

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// BoolString validation errors
var (
	ErrInvalidBoolString = errors.New("invalid BoolString value: expected '0', '1', true, or false")
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

// BoolString is a custom type that can unmarshal boolean JSON values that come as "0"/"1" strings or actual booleans
// Note: UnmarshalJSON uses pointer receiver (must modify value) while MarshalJSON uses value receiver
// (must work with both pointers and values). This is the correct pattern for JSON marshaling.
//
//nolint:recvcheck // Mixed receivers are correct for JSON marshaler/unmarshaler pattern
type BoolString bool

// UnmarshalJSON is a custom unmarshaler for BoolString to handle both string and boolean values
func (bs *BoolString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as boolean first
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		*bs = BoolString(boolValue)
		return nil
	}

	// If that fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return fmt.Errorf("failed to unmarshal BoolString from string: %w", err)
	}

	// Convert string to boolean: "1" = true, "0" = false
	switch stringValue {
	case "1":
		*bs = BoolString(true)
	case "0":
		*bs = BoolString(false)
	default:
		return fmt.Errorf("%w: '%s'", ErrInvalidBoolString, stringValue)
	}

	return nil
}

// MarshalJSON is a custom marshaler for BoolString to output as "0"/"1" strings for API compatibility
func (bs BoolString) MarshalJSON() ([]byte, error) {
	if bs {
		return []byte(`"1"`), nil
	}
	return []byte(`"0"`), nil
}

// Bool converts BoolString back to a standard bool value
func (bs BoolString) Bool() bool {
	return bool(bs)
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

// TagsEvalType represents how tags are evaluated.
type TagsEvalType int

// Tags evaluation types
const (
	// TagsEvalTypeAnd - AND evaluation (all tags must match)
	TagsEvalTypeAnd TagsEvalType = 0
	// TagsEvalTypeOr - OR evaluation (at least one tag must match)
	TagsEvalTypeOr TagsEvalType = 1
)

// UnmarshalJSON is a custom unmarshaler for TagsEvalType to handle both string and int values
func (tet *TagsEvalType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as integer first
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*tet = TagsEvalType(intValue)
		return nil
	}

	// If that fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return fmt.Errorf("failed to unmarshal TagsEvalType from string: %w", err)
	}

	// Convert string to int
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return fmt.Errorf("failed to parse TagsEvalType from string '%s': %w", stringValue, err)
	}

	*tet = TagsEvalType(intValue)
	return nil
}

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
	TimePeriodTypeDaily TimePeriodType = 1
	// TimePeriodTypeWeekly - weekly
	TimePeriodTypeWeekly TimePeriodType = 2
	// TimePeriodTypeMonthly - monthly (by day of month)
	TimePeriodTypeMonthly TimePeriodType = 3
	// TimePeriodTypeMonthlyByWeekday - monthly (by day of week)
	TimePeriodTypeMonthlyByWeekday TimePeriodType = 4
	// TimePeriodTypeYearly - yearly
	TimePeriodTypeYearly TimePeriodType = 5
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
	// Type of tag evaluation (0=AND, 1=OR)
	TagsEvalType TagsEvalType `json:"tags_evaltype,omitempty"`
}

// TimePeriod represents a time period for maintenance.
type TimePeriod struct {
	// Time period ID (readonly)
	TimePeriodID string `json:"timeperiodid,omitempty"`
	// Time period type (0=one time, 1=daily, 2=weekly, 3=monthly by day, 4=monthly by weekday, 5=yearly)
	TimePeriodType TimePeriodType `json:"timeperiod_type"`
	// Start date for one-time maintenance (Unix timestamp) - required for type 0
	StartDate int64 `json:"start_date,omitempty"`
	// For recurring periods, frequency (e.g., every N days/weeks) - required for type 4
	Every int `json:"every,omitempty"`
	// Day of the month (1-31) - required for type 3 and 5
	Day int `json:"day,omitempty"`
	// Day of the week (0=Sunday, 1=Monday, ..., 6=Saturday) - required for type 2 and 4
	DayOfWeek int `json:"dayofweek,omitempty"`
	// Month of the year (1-12) - required for type 5
	Month int `json:"month,omitempty"`
	// Year (optional, for type 5)
	Year int `json:"year,omitempty"`
	// Start time in seconds from midnight - required for types 1,2,3,4,5
	StartTime int `json:"start_time,omitempty"`
	// Duration of the maintenance period in seconds - required for all types
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
	SelectGroups      any      `json:"selectGroups,omitempty"`      // "extend" or array of fields
	SelectHosts       any      `json:"selectHosts,omitempty"`       // "extend" or array of fields
	SelectTags        any      `json:"selectTags,omitempty"`        // "extend" or array of fields
	SelectTimePeriods any      `json:"selectTimeperiods,omitempty"` // "extend" or array of fields
	LimitSelects      int      `json:"limitSelects,omitempty"`      // Limits the number of records returned by subselects
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

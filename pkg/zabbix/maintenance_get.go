package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MaintenanceGetResponse represents the response from maintenance.get API call.
type MaintenanceGetResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  []Maintenance `json:"result,omitempty"`
	Error   Error         `json:"error,omitempty"`
	ID      int           `json:"id"`
}

// MaintenanceGetOption defines a function signature for options to configure a MaintenanceGetRequest.
type MaintenanceGetOption func(*MaintenanceGetRequest)

// NewMaintenanceGetRequest creates a new request for the maintenance.get API method with default values.
func NewMaintenanceGetRequest(options ...MaintenanceGetOption) *MaintenanceGetRequest {
	mgr := &MaintenanceGetRequest{
		JSONRPC: "2.0",
		Method:  "maintenance.get",
		Params:  MaintenanceGetParams{},
		ID:      1, // Default ID
	}
	
	// Apply all options
	for _, opt := range options {
		opt(mgr)
	}
	
	return mgr
}

// WithMaintenanceGetGroupIDs adds host group IDs to filter the returned maintenances.
func WithMaintenanceGetGroupIDs(groupIDs []string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.GroupIDs = groupIDs
	}
}

// WithMaintenanceGetHostIDs adds host IDs to filter the returned maintenances.
func WithMaintenanceGetHostIDs(hostIDs []string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.HostIDs = hostIDs
	}
}

// WithMaintenanceGetMaintenanceIDs adds maintenance IDs to filter the maintenances.
func WithMaintenanceGetMaintenanceIDs(maintenanceIDs []string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.MaintenanceIDs = maintenanceIDs
	}
}

// WithMaintenanceGetOutput sets the output parameter to control which fields are returned.
func WithMaintenanceGetOutput(output interface{}) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.Output = output
	}
}

// WithMaintenanceGetSelectGroups adds the selectGroups parameter to retrieve host groups.
func WithMaintenanceGetSelectGroups(selectGroups string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.SelectGroups = selectGroups
	}
}

// WithMaintenanceGetSelectHosts adds the selectHosts parameter to retrieve hosts.
func WithMaintenanceGetSelectHosts(selectHosts string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.SelectHosts = selectHosts
	}
}

// WithMaintenanceGetSelectTimePeriods adds the selectTimeperiods parameter to retrieve time periods.
func WithMaintenanceGetSelectTimePeriods(selectTimePeriods string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.SelectTimePeriods = selectTimePeriods
	}
}

// WithMaintenanceGetAuthToken sets the authentication token for the API request.
func WithMaintenanceGetAuthToken(token string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Auth = token
	}
}

// WithMaintenanceGetID sets the ID for the API request.
func WithMaintenanceGetID(id int) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.ID = id
	}
}

// WithMaintenanceGetLimit adds a limit to the number of returned maintenances.
func WithMaintenanceGetLimit(limit int) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.Limit = limit
	}
}

// WithMaintenanceGetSortField sets the fields to sort the results by.
func WithMaintenanceGetSortField(sortField []string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.SortField = sortField
	}
}

// WithMaintenanceGetSortOrder sets the sort order of the results.
func WithMaintenanceGetSortOrder(sortOrder []string) MaintenanceGetOption {
	return func(mgr *MaintenanceGetRequest) {
		mgr.Params.SortOrder = sortOrder
	}
}

// MaintenanceGet sends a maintenance.get request to the Zabbix API.
func (z *Client) MaintenanceGet(ctx context.Context, request *MaintenanceGetRequest) (*MaintenanceGetResponse, error) {
	statusCode, respBody, err := z.postRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API request failed for maintenance.get: %w", err)
	}

	if statusCode != http.StatusOK {
		// Try to unmarshal the response body as a Zabbix JSON-RPC error object
		var zbxErrorResponse struct {
			Error *Error `json:"error,omitempty"`
		}
		if unmarshalErr := json.Unmarshal(respBody, &zbxErrorResponse); unmarshalErr == nil && zbxErrorResponse.Error != nil {
			// If we successfully parsed a Zabbix error, return it directly
			return nil, zbxErrorResponse.Error
		}
		// If we couldn't parse a specific Zabbix error, return a generic HTTP error
		return nil, fmt.Errorf("API request for maintenance.get returned HTTP status %d and unexpected response body: %s", statusCode, string(respBody))
	}

	var response MaintenanceGetResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal maintenance.get response: %w - %s", err, string(respBody))
	}

	if response.Error.Code != 0 {
		return nil, response.Error
	}

	return &response, nil
}

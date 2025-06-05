package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MaintenanceDeleteResponse represents the response from the maintenance.delete API call.
type MaintenanceDeleteResponse struct {
	JSONRPC string               `json:"jsonrpc"`
	Result  MaintenanceDeleteResult `json:"result"`
	ID      int                  `json:"id"`
	Error   *Error               `json:"error,omitempty"`
}

// MaintenanceDeleteResult represents the result object in the maintenance.delete response
type MaintenanceDeleteResult struct {
	MaintenanceIDs []string `json:"maintenanceids"` // IDs of deleted maintenance periods
}

// MaintenanceDeleteOption defines a function signature for options to configure a MaintenanceDeleteRequest.
type MaintenanceDeleteOption func(*MaintenanceDeleteRequest)

// NewMaintenanceDeleteRequest creates a new request for deleting maintenance periods.
func NewMaintenanceDeleteRequest(options ...MaintenanceDeleteOption) *MaintenanceDeleteRequest {
	md := &MaintenanceDeleteRequest{
		JSONRPC: JSONRPC,
		Method:  "maintenance.delete",
		Params:  []string{},
		ID:      generateUniqueID(),
	}

	for _, opt := range options {
		opt(md)
	}

	return md
}

// WithMaintenanceDeleteIDs adds maintenance IDs to be deleted.
func WithMaintenanceDeleteIDs(ids []string) MaintenanceDeleteOption {
	return func(md *MaintenanceDeleteRequest) {
		md.Params = ids
	}
}

// WithMaintenanceDeleteID adds a single maintenance ID to be deleted.
func WithMaintenanceDeleteID(id string) MaintenanceDeleteOption {
	return func(md *MaintenanceDeleteRequest) {
		md.Params = append(md.Params, id)
	}
}

// WithMaintenanceDeleteAuthToken sets the authentication token for the request.
func WithMaintenanceDeleteAuthToken(authToken string) MaintenanceDeleteOption {
	return func(md *MaintenanceDeleteRequest) {
		md.Auth = authToken
	}
}

// MaintenanceDelete deletes maintenance periods by their IDs.
// The request object should be fully populated by the caller, including Auth and ID.
func (z *Client) MaintenanceDelete(ctx context.Context, request *MaintenanceDeleteRequest) (*MaintenanceDeleteResponse, error) {
	statusCode, respBody, err := z.postRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API request failed for maintenance.delete: %w", err)
	}

	if statusCode != http.StatusOK {
		// Try to unmarshal the response body as a Zabbix JSON-RPC error object
		var zbxErrorResponse struct {
			Error *Error `json:"error,omitempty"`
		}
		if unmarshalErr := json.Unmarshal(respBody, &zbxErrorResponse); unmarshalErr == nil && 
		   zbxErrorResponse.Error != nil && zbxErrorResponse.Error.Code != 0 {
			// If we successfully parsed a Zabbix error, return it directly
			return nil, zbxErrorResponse.Error
		}
		// If we couldn't parse a specific Zabbix error, return a generic HTTP error
		return nil, fmt.Errorf("API request for maintenance.delete returned HTTP status %d and unexpected response body: %s", 
			statusCode, string(respBody))
	}

	var response MaintenanceDeleteResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal maintenance.delete response: %w - %s", err, string(respBody))
	}

	if response.Error != nil && response.Error.Code != 0 {
		return nil, response.Error // response.Error already implements the error interface
	}

	return &response, nil
}

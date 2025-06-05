package zabbix

import (
	"context"
	"fmt"
)

// MaintenanceCreateResponse represents the response from maintenance.create API call.
type MaintenanceCreateResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  MaintenanceResponse `json:"result,omitempty"`
	Error   *Error     `json:"error,omitempty"`
	ID      int              `json:"id"`
}

// MaintenanceCreate sends a maintenance.create request to the Zabbix API.
// The request object should be fully populated by the caller, including Auth and ID.
func (z *Client) MaintenanceCreate(ctx context.Context, request *MaintenanceCreateRequest) (*MaintenanceCreateResponse, error) {
	statusCode, respBody, err := z.postRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API request failed for maintenance.create: %w", err)
	}

	var response MaintenanceCreateResponse
	if err := handleRawResponse(statusCode, respBody, "maintenance.create", &response); err != nil {
		return nil, err
	}

	// After handleRawResponse, we need to check the Error field of the specific response type
	if response.Error != nil && response.Error.Code != 0 {
		return nil, response.Error
	}

	return &response, nil
}



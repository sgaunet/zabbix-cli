package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MaintenanceCreateResponse represents the response from maintenance.create API call.
type MaintenanceCreateResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  MaintenanceResponse `json:"result,omitempty"`
	Error   Error      `json:"error,omitempty"` // Using standard Error struct as required by memory
	ID      int              `json:"id"`
}

// MaintenanceCreate sends a maintenance.create request to the Zabbix API.
// The request object should be fully populated by the caller, including Auth and ID.
func (z *Client) MaintenanceCreate(ctx context.Context, request *MaintenanceCreateRequest) (*MaintenanceCreateResponse, error) {
	statusCode, respBody, err := z.postRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API request failed for maintenance.create: %w", err)
	}

	if statusCode != http.StatusOK {
		// Try to unmarshal the response body as a Zabbix JSON-RPC error object.
		var zbxErrorResponse struct {
			Error *Error `json:"error,omitempty"`
		}
		if unmarshalErr := json.Unmarshal(respBody, &zbxErrorResponse); unmarshalErr == nil && zbxErrorResponse.Error != nil {
			// If we successfully parsed a Zabbix error, return it directly.
			return nil, zbxErrorResponse.Error
		}
		// If we couldn't parse a specific Zabbix error, return a generic HTTP error.
		return nil, fmt.Errorf("API request for maintenance.create returned HTTP status %d and unexpected response body: %s", statusCode, string(respBody))
	}

	var response MaintenanceCreateResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal maintenance.create response: %w - %s", err, string(respBody))
	}

	if response.Error.Code != 0 {
		return nil, response.Error
	}

	return &response, nil
}



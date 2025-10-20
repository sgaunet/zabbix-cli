package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Documentation of zabbix api: https://www.zabbix.com/documentation/7.2/en/manual/api/reference/configuration/import

const methodConfigurationImport = "configuration.import"

// configurationImportRequest is the request body for configuration.import.
// it's the same as configuration.importcompare.
type configurationImportRequest configurationImportCompareRequest

// NewConfigurationImportRequest creates a new configuration import request with default rules.
func NewConfigurationImportRequest(source string) *configurationImportRequest {
	return &configurationImportRequest{
		JSONRPC: JSONRPC,
		Method:  methodConfigurationImport,
		Params: paramsImport{
			Source: source,
			Format: "yaml",
			Rules:  rulesAllTrue(),
		},
		Auth: "",
		ID:   generateUniqueID(),
	}
}

type configurationImportResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  bool   `json:"result"`
	Error   *Error `json:"error,omitempty"`
	ID      int    `json:"id"`
}

// Import imports configuration from the given source string.
func (z *Client) Import(ctx context.Context, source string) (bool, error) {
	c := NewConfigurationImportRequest(source)
	// initialize auth token
	c.Auth = z.Auth()

	statusCode, body, err := z.postRequest(ctx, c)
	if err != nil {
		return false, fmt.Errorf("cannot do request: %w", err)
	}

	if statusCode != http.StatusOK {
		return false, fmt.Errorf("status code not OK: %d - %s (%w)", statusCode, string(body), ErrWrongHTTPCode)
	}

	var res configurationImportResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, fmt.Errorf("cannot unmarshal response: %w - %s", err, string(body))
	}
	if res.Error != nil && res.Error.Code != 0 {
		return false, res.Error
	}
	return res.Result, nil
}

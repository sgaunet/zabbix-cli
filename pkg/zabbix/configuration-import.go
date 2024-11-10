package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Documentation of zabbix api: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/configuration/import

const methodConfigurationImport = "configuration.import"

// configurationImportRequest is the request body for configuration.import
// it's the same as configuration.importcompare
type configurationImportRequest configurationImportCompareRequest

func newConfigurationImportRequest(source string) *configurationImportRequest {
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
	JSONRPC  string   `json:"jsonrpc"`
	Result   bool     `json:"result"`
	ErrorMsg ErrorMsg `json:"error,omitempty"`
	ID       int      `json:"id"`
}

func (z *ZabbixAPI) Import(ctx context.Context, source string) (bool, error) {
	c := newConfigurationImportRequest(source)
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
	if res.ErrorMsg != (ErrorMsg{}) {
		return false, fmt.Errorf("error message: %w", &res.ErrorMsg)
	}
	return res.Result, nil
}

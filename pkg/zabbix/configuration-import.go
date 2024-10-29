package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	postBody, err := json.Marshal(c)
	if err != nil {
		return false, fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, z.APIEndpoint, responseBody)
	if err != nil {
		return false, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("cannot read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("status code not OK: %d - %s", resp.StatusCode, string(body))
	}

	var res configurationImportResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, fmt.Errorf("cannot unmarshal response: %w - %s", err, string(body))
	}
	if res.ErrorMsg != (ErrorMsg{}) {
		return false, fmt.Errorf("error message: %s", res.ErrorMsg.Error())
	}
	return res.Result, nil
}

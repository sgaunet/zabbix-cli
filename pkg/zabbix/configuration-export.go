package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Documentation of zabbix api: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/configuration/export

const MethodConfigurationExport = "configuration.export"

type ConfigurationExportRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Options struct {
			Templates []string `json:"templates"`
		} `json:"options"`
		Format string `json:"format"`
	} `json:"params"`
	Auth string `json:"auth"`
	ID   int    `json:"id"`
}

type ConfigurationExportResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
}

func NewConfigurationExportRequest(templates []string) *ConfigurationExportRequest {
	return &ConfigurationExportRequest{
		JSONRPC: JSONRPC,
		Method:  MethodConfigurationExport,
		Params: struct {
			Options struct {
				Templates []string `json:"templates"`
			} `json:"options"`
			Format string `json:"format"`
		}{
			Options: struct {
				Templates []string `json:"templates"`
			}{
				Templates: templates,
			},
			Format: "yaml",
		},
		Auth: "",
		ID:   generateUniqueID(),
	}
}

func (z *ZabbixAPI) Export(ctx context.Context, c *ConfigurationExportRequest) (*ConfigurationExportResponse, error) {
	// initialize auth token
	c.Auth = z.Auth()
	postBody, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, z.APIEndpoint, responseBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}

	var res ConfigurationExportResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return &res, nil
}

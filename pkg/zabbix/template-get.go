package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const MethodTemplateGet = "template.get"

// TemplateGetRequest struct is used to get templates from the Zabbix API
type TemplateGetRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Output string       `json:"output"`
		Filter ZabbixFilter `json:"filter"`
	} `json:"params"`
	Auth string `json:"auth"`
	ID   int    `json:"id"`
}

// TemplateGetResponse struct is used to unmarshal the response from the Zabbix API
type TemplateGetResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  []struct {
		TemplateID string `json:"templateID"`
		Name       string `json:"name"`
	} `json:"result"`
	ID int `json:"id"`
}

// NewTemplateGetRequest returns a new TemplateGetRequest
func (z *ZabbixAPI) NewTemplateGetRequest(templatesNames []string) *TemplateGetRequest {
	return &TemplateGetRequest{
		JSONRPC: JSONRPC,
		Method:  MethodTemplateGet,
		Params: struct {
			Output string       `json:"output"`
			Filter ZabbixFilter `json:"filter"`
		}{
			Output: "extend",
			Filter: ZabbixFilter{
				Name: templatesNames,
			},
		},
		Auth: z.Auth(),
		ID:   1,
	}
}

// GetTemplate returns the templates from the Zabbix API
func (z *ZabbixAPI) GetTemplate(c *TemplateGetRequest) (*TemplateGetResponse, error) {
	// initialize auth token
	c.Auth = z.Auth()
	postBody, _ := json.Marshal(c)
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, z.APIEndpoint, responseBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(context.TODO())
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

	var res TemplateGetResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return &res, nil
}

// GetTemplateID returns the ID of the templates
func (t *TemplateGetResponse) GetTemplateID() []string {
	var templates []string
	for _, template := range t.Result {
		templates = append(templates, template.TemplateID)
	}
	return templates
}

// GetTemplateName returns the name of the templates
func (t *TemplateGetResponse) GetTemplateName() []string {
	var templates []string
	for _, template := range t.Result {
		templates = append(templates, template.Name)
	}
	return templates
}

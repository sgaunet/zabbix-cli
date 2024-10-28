package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Documentation of zabbix api: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/template/get

const methodTemplateGet = "template.get"

type templateGetRequestParams struct {
	Output string                 `json:"output"`
	Filter map[string]interface{} `json:"filter"`
}

// templateGetRequest struct is used to get templates from the Zabbix API
type templateGetRequest struct {
	JSONRPC string                   `json:"jsonrpc"`
	Method  string                   `json:"method"`
	Params  templateGetRequestParams `json:"params"`
	Auth    string                   `json:"auth"`
	ID      int                      `json:"id"`
}

// templateGetResponse struct is used to unmarshal the response from the Zabbix API
type templateGetResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  []struct {
		TemplateID string `json:"templateID"`
		Name       string `json:"name"`
	} `json:"result"`
	ID int `json:"id"`
}

// NewTemplateGetRequest returns a new TemplateGetRequest
func (z *ZabbixAPI) GetTemplates(templatesNames []string) (*templateGetResponse, error) {
	payload := &templateGetRequest{
		JSONRPC: JSONRPC,
		Method:  methodTemplateGet,
		Params: templateGetRequestParams{
			Output: "extend",
			Filter: NewZabbixFilterGetMethod(AddFilter("name", templatesNames)).GetFilter(),
		},
		Auth: z.Auth(),
		ID:   generateUniqueID(),
	}

	postBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal data: %w", err)
	}
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

	var res templateGetResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return &res, nil
}

// // GetTemplate returns the templates from the Zabbix API
// func (z *ZabbixAPI) GetTemplate(c *templateGetRequest) (*templateGetResponse, error) {
// 	// initialize auth token
// 	c.Auth = z.Auth()
// 	postBody, _ := json.Marshal(c)
// 	responseBody := bytes.NewBuffer(postBody)
// 	req, err := http.NewRequest(http.MethodPost, z.APIEndpoint, responseBody)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create request: %w", err)
// 	}
// 	req = req.WithContext(context.TODO())
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := z.client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot do request: %w", err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot read response body: %w", err)
// 	}

// 	var res templateGetResponse
// 	err = json.Unmarshal(body, &res)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
// 	}
// 	return &res, nil
// }

// GetTemplateID returns the ID of the templates
func (t *templateGetResponse) GetTemplateID() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.TemplateID)
	}
	return templates
}

// GetTemplateName returns the name of the templates
func (t *templateGetResponse) GetTemplateName() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.Name)
	}
	return templates
}

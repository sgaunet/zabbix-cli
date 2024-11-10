package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
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

// templateGetRequestOption is the options struct to get templates
type templateGetRequestOption func(*templateGetRequest)

// GetTemplateOptionFilterByName returns a templateGetRequestOption to filter by name
func GetTemplateOptionFilterByName(templatesNames []string) templateGetRequestOption {
	return func(c *templateGetRequest) {
		c.Params.Filter["name"] = templatesNames
	}
}

// newTemplateGetRequest returns a new templateGetRequest
func newTemplateGetRequest(options ...templateGetRequestOption) *templateGetRequest {
	c := &templateGetRequest{
		JSONRPC: JSONRPC,
		Method:  methodTemplateGet,
		Params: templateGetRequestParams{
			Output: "extend",
			Filter: map[string]interface{}{},
		},
		Auth: "",
		ID:   generateUniqueID(),
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

// templateGetResponse struct is used to unmarshal the response from the Zabbix API
type templateGetResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  []struct {
		TemplateID string `json:"templateID"`
		Name       string `json:"name"`
	} `json:"result"`
	// Error is the error message when the request fails TODO
	ID int `json:"id"`
}

// NewTemplateGetRequest returns a new TemplateGetRequest
func (z *ZabbixAPI) GetTemplates(options ...templateGetRequestOption) (*templateGetResponse, error) {
	payload := newTemplateGetRequest(options...)
	payload.Auth = z.Auth()

	statusCode, body, err := z.postRequest(context.Background(), payload)
	if err != nil {
		return nil, fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%w)", statusCode, ErrWrongHTTPCode)
	}

	var res templateGetResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	return &res, nil
}

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

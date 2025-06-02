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

// TemplateGetRequestOption is the options struct to get templates
// TemplateGetRequestOption is a function that modifies a templateGetRequest.
type TemplateGetRequestOption func(*templateGetRequest)

// GetTemplateOptionFilterByName returns a TemplateGetRequestOption to filter by name
// GetTemplateOptionFilterByName returns a TemplateGetRequestOption to filter by name.
func GetTemplateOptionFilterByName(templatesNames []string) TemplateGetRequestOption {
	return func(c *templateGetRequest) {
		c.Params.Filter["name"] = templatesNames
	}
}

// newTemplateGetRequest returns a new templateGetRequest
func newTemplateGetRequest(options ...TemplateGetRequestOption) *templateGetRequest {
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

// TemplateGetResponse struct is used to unmarshal the response from the Zabbix API
type TemplateGetResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  []struct {
		TemplateID string `json:"templateID"`
		Name       string `json:"name"`
	} `json:"result"`
	ErrorMsg ErrorMsg `json:"error,omitempty"`
	ID       int      `json:"id"`
}

// GetTemplates returns templates from the Zabbix API.
// GetTemplates returns templates from the Zabbix API.
// GetTemplates returns templates from the Zabbix API. Note: The return type TemplateGetResponse is unexported by design for internal encapsulation.
func (z *Client) GetTemplates(options ...TemplateGetRequestOption) (*TemplateGetResponse, error) {
	payload := newTemplateGetRequest(options...)
	payload.Auth = z.Auth()

	statusCode, body, err := z.postRequest(context.Background(), payload)
	if err != nil {
		return nil, fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%w)", statusCode, ErrWrongHTTPCode)
	}

	var res TemplateGetResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	if res.ErrorMsg != (ErrorMsg{}) {
		return nil, fmt.Errorf("error message: %w", &res.ErrorMsg)
	}
	return &res, nil
}

// GetTemplateID returns the ID of the templates
func (t *TemplateGetResponse) GetTemplateID() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.TemplateID)
	}
	return templates
}

// GetTemplateName returns the name of the templates
func (t *TemplateGetResponse) GetTemplateName() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.Name)
	}
	return templates
}

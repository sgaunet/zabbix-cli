package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Documentation of zabbix api: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/configuration/export
const methodConfigurationExport = "configuration.export"
const defaultConfigurationExportFormat = "yaml"

// ConfigurationExportOptions is the options struct to export configuration.
type ConfigurationExportOptions struct {
	HostGroupsID     []string `json:"host_groups,omitempty"`     // (array) IDs of host groups to export (Zabbix 6.2+).
	TemplateGroupsID []string `json:"template_groups,omitempty"` // (array) IDs of template groups to export (Zabbix 6.2+).
	HostsID          []string `json:"hosts,omitempty"`           // (array) IDs of hosts to export.
	ImagesID         []string `json:"images,omitempty"`          // (array) IDs of images to export.
	MapsID           []string `json:"maps,omitempty"`            // (array) IDs of maps to export.
	MediaTypesID     []string `json:"mediaTypes,omitempty"`      // (array) IDs of media types to export.
	TemplatesID      []string `json:"templates,omitempty"`       // (array) IDs of templates to export.
	DashboardsID     []string `json:"dashboards,omitempty"`      // (array) IDs of dashboards to export.
}

// ConfigurationExportParams is the params struct to export configuration.
type ConfigurationExportParams struct {
	Options ConfigurationExportOptions `json:"options"`
	Format  string                     `json:"format"`
}

// ConfigurationExportRequest is the request struct to export configuration.
type ConfigurationExportRequest struct {
	JSONRPC string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  ConfigurationExportParams `json:"params"`
	Auth    string                    `json:"auth"`
	ID      int                       `json:"id"`
}

// ConfigurationExportRequestOption is a function that modifies a ConfigurationExportRequest.
type ConfigurationExportRequestOption func(*ConfigurationExportRequest)

// NewConfigurationExportRequest creates a new configuration export request. Note: The return type ConfigurationExportRequest is unexported by design for internal encapsulation.
func NewConfigurationExportRequest(options ...ConfigurationExportRequestOption) *ConfigurationExportRequest {
	c := &ConfigurationExportRequest{
		JSONRPC: JSONRPC,
		Method:  methodConfigurationExport,
		Params: ConfigurationExportParams{
			Options: ConfigurationExportOptions{},
			Format:  defaultConfigurationExportFormat,
		},
		Auth: "",
		ID:   generateUniqueID(),
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

// ExportRequestOptionHostGroupsID returns a ConfigurationExportRequestOption for host groups.
func ExportRequestOptionHostGroupsID(groupsID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.HostGroupsID = groupsID
	}
}

// ExportRequestOptionTemplateGroupsID returns a ConfigurationExportRequestOption for template groups.
func ExportRequestOptionTemplateGroupsID(groupsID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.TemplateGroupsID = groupsID
	}
}

// ExportRequestOptionHostsID returns a ConfigurationExportRequestOption for hosts.
func ExportRequestOptionHostsID(hostsID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.HostsID = hostsID
	}
}

// ExportRequestOptionImagesID returns a ConfigurationExportRequestOption for images.
func ExportRequestOptionImagesID(imagesID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.ImagesID = imagesID
	}
}

// ExportRequestOptionMapsID returns a ConfigurationExportRequestOption for maps.
func ExportRequestOptionMapsID(mapsID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.MapsID = mapsID
	}
}

// ExportRequestOptionMediaTypesID returns a ConfigurationExportRequestOption for media types.
func ExportRequestOptionMediaTypesID(mediaTypesID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.MediaTypesID = mediaTypesID
	}
}

// ExportRequestOptionTemplatesID returns a ConfigurationExportRequestOption for templates.
func ExportRequestOptionTemplatesID(templatesID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.TemplatesID = templatesID
	}
}

// ExportRequestOptionYAMLFormat returns a ConfigurationExportRequestOption for YAML format.
// ExportRequestOptionYAMLFormat returns a ConfigurationExportRequestOption for YAML format.
func ExportRequestOptionYAMLFormat() ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Format = "yaml"
	}
}

// ExportRequestOptionJSONFormat returns a ConfigurationExportRequestOption for JSON format.
// ExportRequestOptionJSONFormat returns a ConfigurationExportRequestOption for JSON format.
func ExportRequestOptionJSONFormat() ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Format = "json"
	}
}

// ExportRequestOptionXMLFormat returns a ConfigurationExportRequestOption for XML format.
// ExportRequestOptionXMLFormat returns a ConfigurationExportRequestOption for XML format.
func ExportRequestOptionXMLFormat() ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Format = "xml"
	}
}

// ExportRequestOptionDashboardsID returns a ConfigurationExportRequestOption for dashboards.
func ExportRequestOptionDashboardsID(dashboardsID []string) ConfigurationExportRequestOption {
	return func(c *ConfigurationExportRequest) {
		c.Params.Options.DashboardsID = dashboardsID
	}
}

// ConfigurationExportResponse is the response struct of a configuration export request.
type ConfigurationExportResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      int    `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

// Export exports the configuration using the provided options.
func (z *Client) Export(ctx context.Context, opt ...ConfigurationExportRequestOption) (string, error) {
	c := NewConfigurationExportRequest(opt...)
	// initialize auth token
	c.Auth = z.Auth()
	statusCode, body, err := z.postRequest(ctx, c)
	if err != nil {
		return "", fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d (%w)", statusCode, ErrWrongHTTPCode)
	}

	var res ConfigurationExportResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal response: %w", err)
	}
	if res.Error != nil && res.Error.Code != 0 {
		return "", res.Error
	}
	return res.Result, nil
}

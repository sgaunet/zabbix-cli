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

// configurationExportOptions is the options struct to export configuration
type configurationExportOptions struct {
	GroupsID     []string `json:"groups,omitempty"`     // (array) IDs of host groups to export.
	HostsID      []string `json:"hosts,omitempty"`      // (array) IDs of hosts to export.
	ImagesID     []string `json:"images,omitempty"`     // (array) IDs of images to export.
	MapsID       []string `json:"maps,omitempty"`       // (array) IDs of maps to export.
	MediaTypesID []string `json:"mediaTypes,omitempty"` // (array) IDs of media types to export.
	TemplatesID  []string `json:"templates,omitempty"`  // (array) IDs of templates to export.
}

// configurationExportParams is the params struct to export configuration
type configurationExportParams struct {
	Options configurationExportOptions `json:"options"`
	Format  string                     `json:"format"`
}

// configurationExportRequest is the request struct to export configuration
type configurationExportRequest struct {
	JSONRPC string                    `json:"jsonrpc"`
	Method  string                    `json:"method"`
	Params  configurationExportParams `json:"params"`
	Auth    string                    `json:"auth"`
	ID      int                       `json:"id"`
}

type configurationExportRequestOption func(*configurationExportRequest)

func NewConfigurationExportRequest(options ...configurationExportRequestOption) *configurationExportRequest {
	c := &configurationExportRequest{
		JSONRPC: JSONRPC,
		Method:  methodConfigurationExport,
		Params: configurationExportParams{
			Options: configurationExportOptions{},
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

func ExportRequestOptionGroupsID(groupsID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.GroupsID = groupsID
	}
}

func ExportRequestOptionHostsID(hostsID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.HostsID = hostsID
	}
}

func ExportRequestOptionImagesID(imagesID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.ImagesID = imagesID
	}
}

func ExportRequestOptionMapsID(mapsID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.MapsID = mapsID
	}
}

func ExportRequestOptionMediaTypesID(mediaTypesID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.MediaTypesID = mediaTypesID
	}
}

func ExportRequestOptionTemplatesID(templatesID []string) configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Options.TemplatesID = templatesID
	}
}

func ExportRequestOptionYAMLFormat() configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Format = "yaml"
	}
}

func ExportRequestOptionJSONFormat() configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Format = "json"
	}
}

func ExportRequestOptionXMLFormat() configurationExportRequestOption {
	return func(c *configurationExportRequest) {
		c.Params.Format = "xml"
	}
}

// ConfigurationExportRequest is the response struct of a configuration export request
type ConfigurationExportResponse struct {
	JSONRPC  string   `json:"jsonrpc"`
	Result   string   `json:"result"`
	ID       int      `json:"id"`
	ErrorMsg ErrorMsg `json:"error,omitempty"`
}

func (z *ZabbixAPI) Export(ctx context.Context, opt ...configurationExportRequestOption) (string, error) {
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
	if res.ErrorMsg != (ErrorMsg{}) {
		return "", fmt.Errorf("error message: %w", &res.ErrorMsg)
	}
	return res.Result, nil
}

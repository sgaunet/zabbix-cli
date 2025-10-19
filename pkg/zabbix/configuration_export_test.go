package zabbix_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationExport(t *testing.T) {
	t.Parallel()

	t.Run("NewConfigurationExportRequest", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest()
		require.NotNil(t, c)
		require.Equal(t, zabbix.JSONRPC, c.JSONRPC)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
	})

	t.Run("ExportRequestOptionHostGroupsID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionHostGroupsID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.HostGroupsID)
	})

	t.Run("ExportRequestOptionTemplateGroupsID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionTemplateGroupsID([]string{"3", "4"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"3", "4"}, c.Params.Options.TemplateGroupsID)
	})

	t.Run("ExportRequestOptionHostsID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionHostsID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.HostsID)
	})

	t.Run("ExportRequestOptionImagesID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionImagesID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.ImagesID)
	})

	t.Run("ExportRequestOptionMapsID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionMapsID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.MapsID)
	})

	t.Run("ExportRequestOptionMediaTypesID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionMediaTypesID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.MediaTypesID)
	})

	t.Run("ExportRequestOptionTemplatesID", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionTemplatesID([]string{"1", "2"}))
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.TemplatesID)
	})

	t.Run("ExportRequestOptionYAMLFormat", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionYAMLFormat())
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, "yaml", c.Params.Format)
	})

	t.Run("ExportRequestOptionJSONFormat", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionJSONFormat())
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, "json", c.Params.Format)
	})

	t.Run("ExportRequestOptionXMLFormat", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionXMLFormat())
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, "xml", c.Params.Format)
	})

	t.Run("ExportRequestOptionAll", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(
			zabbix.ExportRequestOptionHostGroupsID([]string{"1", "2"}),
			zabbix.ExportRequestOptionTemplateGroupsID([]string{"3", "4"}),
			zabbix.ExportRequestOptionHostsID([]string{"1", "2"}),
			zabbix.ExportRequestOptionImagesID([]string{"1", "2"}),
			zabbix.ExportRequestOptionMapsID([]string{"1", "2"}),
			zabbix.ExportRequestOptionMediaTypesID([]string{"1", "2"}),
			zabbix.ExportRequestOptionTemplatesID([]string{"1", "2"}),
			zabbix.ExportRequestOptionDashboardsID([]string{"5", "6"}),
			zabbix.ExportRequestOptionYAMLFormat(),
		)
		require.NotNil(t, c)
		require.Equal(t, "", c.Auth)
		require.NotEqual(t, 0, c.ID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.HostGroupsID)
		require.Equal(t, []string{"3", "4"}, c.Params.Options.TemplateGroupsID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.HostsID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.ImagesID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.MapsID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.MediaTypesID)
		require.Equal(t, []string{"1", "2"}, c.Params.Options.TemplatesID)
		require.Equal(t, []string{"5", "6"}, c.Params.Options.DashboardsID)
		require.Equal(t, "yaml", c.Params.Format)
	})
}

func TestConfigurationExportRequestMarshal(t *testing.T) {
	t.Parallel()

	t.Run("Test that payload contains only selected options", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(zabbix.ExportRequestOptionTemplatesID([]string{"1", "2"}))
		payload, err := json.Marshal(c)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		require.Contains(t, string(payload), "templates")
		require.NotContains(t, string(payload), "host_groups")
		require.NotContains(t, string(payload), "template_groups")
		require.NotContains(t, string(payload), "hosts")
		require.NotContains(t, string(payload), "images")
		require.NotContains(t, string(payload), "maps")
		require.NotContains(t, string(payload), "mediaTypes")
	})

	t.Run("An empty option should not be taken into account", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(
			zabbix.ExportRequestOptionTemplatesID([]string{"1", "2"}),
			zabbix.ExportRequestOptionHostGroupsID([]string{}),
		)
		payload, err := json.Marshal(c)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		require.Contains(t, string(payload), "templates")
		require.NotContains(t, string(payload), "host_groups")
		require.NotContains(t, string(payload), "template_groups")
		require.NotContains(t, string(payload), "hosts")
		require.NotContains(t, string(payload), "images")
		require.NotContains(t, string(payload), "maps")
		require.NotContains(t, string(payload), "mediaTypes")
	})
}

func TestExport(t *testing.T) {
	t.Parallel()
	t.Run("Export success", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "exported_configuration",
					ID:      1,
				}
				respJSON, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(respJSON) //nolint: errcheck
			}))
		client := ts.Client()
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)
		z.SetHTTPClient(client)
		// No need to login
		resp, err := z.Export(context.Background())
		require.NoError(t, err)
		require.Equal(t, "exported_configuration", resp)
	})

	t.Run("Export failure", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "",
					ID:      1,
				}
				respJSON, err := json.Marshal(resp)
				if err != nil {
					t.Fatal(err)
				}
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(respJSON) //nolint:errcheck
			}))
		client := ts.Client()
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)
		z.SetHTTPClient(client)
		// No need to login
		resp, err := z.Export(context.Background())
		require.Error(t, err)
		require.Empty(t, resp)
	})
}

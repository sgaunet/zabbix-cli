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
		require.NotContains(t, string(payload), "dashboards")
	})

	t.Run("Dashboards option marshals correctly", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(
			zabbix.ExportRequestOptionDashboardsID([]string{"10", "20"}),
			zabbix.ExportRequestOptionJSONFormat(),
		)
		payload, err := json.Marshal(c)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		require.Contains(t, string(payload), "dashboards")
		require.Contains(t, string(payload), `"dashboards":["10","20"]`)
		require.Contains(t, string(payload), `"format":"json"`)
	})

	t.Run("Multiple option types combined", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(
			zabbix.ExportRequestOptionHostsID([]string{"1"}),
			zabbix.ExportRequestOptionTemplatesID([]string{"2"}),
			zabbix.ExportRequestOptionMapsID([]string{"3"}),
			zabbix.ExportRequestOptionXMLFormat(),
		)
		payload, err := json.Marshal(c)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		require.Contains(t, string(payload), "hosts")
		require.Contains(t, string(payload), "templates")
		require.Contains(t, string(payload), "maps")
		require.Contains(t, string(payload), `"format":"xml"`)
		require.NotContains(t, string(payload), "images")
		require.NotContains(t, string(payload), "dashboards")
	})

	t.Run("All options produce valid JSON structure", func(t *testing.T) {
		t.Parallel()
		c := zabbix.NewConfigurationExportRequest(
			zabbix.ExportRequestOptionHostGroupsID([]string{"1"}),
			zabbix.ExportRequestOptionTemplateGroupsID([]string{"2"}),
			zabbix.ExportRequestOptionHostsID([]string{"3"}),
			zabbix.ExportRequestOptionImagesID([]string{"4"}),
			zabbix.ExportRequestOptionMapsID([]string{"5"}),
			zabbix.ExportRequestOptionMediaTypesID([]string{"6"}),
			zabbix.ExportRequestOptionTemplatesID([]string{"7"}),
			zabbix.ExportRequestOptionDashboardsID([]string{"8"}),
		)
		payload, err := json.Marshal(c)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		// Verify we can unmarshal it back
		var req zabbix.ConfigurationExportRequest
		err = json.Unmarshal(payload, &req)
		require.NoError(t, err)
		require.Equal(t, []string{"1"}, req.Params.Options.HostGroupsID)
		require.Equal(t, []string{"2"}, req.Params.Options.TemplateGroupsID)
		require.Equal(t, []string{"3"}, req.Params.Options.HostsID)
		require.Equal(t, []string{"4"}, req.Params.Options.ImagesID)
		require.Equal(t, []string{"5"}, req.Params.Options.MapsID)
		require.Equal(t, []string{"6"}, req.Params.Options.MediaTypesID)
		require.Equal(t, []string{"7"}, req.Params.Options.TemplatesID)
		require.Equal(t, []string{"8"}, req.Params.Options.DashboardsID)
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

	t.Run("Export with YAML format", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request contains correct format
				var req zabbix.ConfigurationExportRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				require.NoError(t, err)
				require.Equal(t, "yaml", req.Params.Format)

				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "zabbix_export:\n  version: '7.2'\n  templates:\n    - template: Test",
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionYAMLFormat(), zabbix.ExportRequestOptionTemplatesID([]string{"1"}))
		require.NoError(t, err)
		require.Contains(t, resp, "zabbix_export:")
	})

	t.Run("Export with JSON format", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request contains correct format
				var req zabbix.ConfigurationExportRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				require.NoError(t, err)
				require.Equal(t, "json", req.Params.Format)

				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  `{"zabbix_export":{"version":"7.2","templates":[{"template":"Test"}]}}`,
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionJSONFormat(), zabbix.ExportRequestOptionTemplatesID([]string{"1"}))
		require.NoError(t, err)
		require.Contains(t, resp, "zabbix_export")
		// Verify it's valid JSON
		var result map[string]interface{}
		err = json.Unmarshal([]byte(resp), &result)
		require.NoError(t, err)
	})

	t.Run("Export with XML format", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request contains correct format
				var req zabbix.ConfigurationExportRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				require.NoError(t, err)
				require.Equal(t, "xml", req.Params.Format)

				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  `<?xml version="1.0" encoding="UTF-8"?><zabbix_export><version>7.2</version><templates><template><template>Test</template></template></templates></zabbix_export>`,
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionXMLFormat(), zabbix.ExportRequestOptionTemplatesID([]string{"1"}))
		require.NoError(t, err)
		require.Contains(t, resp, "<?xml version")
		require.Contains(t, resp, "<zabbix_export>")
	})

	t.Run("Export API error - missing object", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "",
					ID:      1,
					Error: &zabbix.Error{
						Code:    -32500,
						Message: "Application error.",
						Data:    "No permissions to referred object or it does not exist!",
					},
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionTemplatesID([]string{"999999"}))
		require.Error(t, err)
		require.Empty(t, resp)
		require.Contains(t, err.Error(), "No permissions to referred object or it does not exist")
	})

	t.Run("Export API error - permission denied", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "",
					ID:      1,
					Error: &zabbix.Error{
						Code:    -32600,
						Message: "No permissions to access API method",
						Data:    "User does not have permission to call this method.",
					},
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionTemplatesID([]string{"1"}))
		require.Error(t, err)
		require.Empty(t, resp)
		require.Contains(t, err.Error(), "No permissions")
	})

	t.Run("Export malformed JSON response", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"invalid json`)) //nolint: errcheck
			}))
		client := ts.Client()
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)
		z.SetHTTPClient(client)
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionTemplatesID([]string{"1"}))
		require.Error(t, err)
		require.Empty(t, resp)
		require.Contains(t, err.Error(), "cannot unmarshal response")
	})

	t.Run("Export response returns string correctly", func(t *testing.T) {
		t.Parallel()
		expectedYAML := `zabbix_export:
  version: '7.2'
  templates:
    - uuid: test-uuid
      template: Linux by Zabbix agent
      name: Linux by Zabbix agent
      groups:
        - name: Templates/Operating systems`

		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := zabbix.ConfigurationExportResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  expectedYAML,
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
		resp, err := z.Export(context.Background(), zabbix.ExportRequestOptionTemplatesID([]string{"10001"}))
		require.NoError(t, err)
		require.Equal(t, expectedYAML, resp)
		require.IsType(t, "", resp) // Verify it's a string
	})
}

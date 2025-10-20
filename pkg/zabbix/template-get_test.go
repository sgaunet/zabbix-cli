package zabbix_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestGetTemplateID(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := fmt.Sprintf(`{
				"jsonrpc": "%s",
				"result": [
					{
						"templateid": "10001",
						"name": "Template OS Linux"
					}
				],
				"id": 1
			}`, zabbix.JSONRPC)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(res))
		}))
	defer ts.Close()
	// get request
	client := ts.Client()

	z := zabbix.New("user", "password", ts.URL)
	z.SetHTTPClient(client)
	pb, err := z.GetProblems(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, pb)
}

func TestGetTemplateID_Internal(t *testing.T) {
	resp := &zabbix.TemplateGetResponse{
		Result: []zabbix.Template{
			{TemplateID: "10001", Name: "Template OS Linux"},
			{TemplateID: "10002", Name: "Template OS Windows"},
		},
	}
	ids := resp.GetTemplateID()
	require.Equal(t, []string{"10001", "10002"}, ids)
}

func TestGetTemplateName_Internal(t *testing.T) {
	resp := &zabbix.TemplateGetResponse{
		Result: []zabbix.Template{
			{TemplateID: "10001", Name: "Template OS Linux"},
			{TemplateID: "10002", Name: "Template OS Windows"},
		},
	}
	names := resp.GetTemplateName()
	require.Equal(t, []string{"Template OS Linux", "Template OS Windows"}, names)
}

// TestTemplateGetBasic tests basic template.get response structure with all fields.
func TestTemplateGetBasic(t *testing.T) {
	// Verify complete Template structure with all Zabbix API 7.2 fields
	template := zabbix.Template{
		TemplateID:  "10001",
		Host:        "linux-template",
		Name:        "Template OS Linux",
		Description: "Linux OS monitoring template",
		UUID:        "abc123def456",
	}
	template.Vendor.Name = "Zabbix"
	template.Vendor.Version = "7.2"

	// Verify all fields are strings and vendor struct is properly accessible
	require.Equal(t, "10001", template.TemplateID)
	require.Equal(t, "linux-template", template.Host)
	require.Equal(t, "Template OS Linux", template.Name)
	require.Equal(t, "Linux OS monitoring template", template.Description)
	require.Equal(t, "abc123def456", template.UUID)
	require.Equal(t, "Zabbix", template.Vendor.Name)
	require.Equal(t, "7.2", template.Vendor.Version)
}

// TestTemplateGetResponseUnmarshaling tests unmarshaling of template.get JSON response.
func TestTemplateGetResponseUnmarshaling(t *testing.T) {
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"templateid": "10001",
				"host": "linux-template",
				"name": "Template OS Linux",
				"description": "Linux OS monitoring",
				"uuid": "abc123",
				"vendor": {
					"name": "Zabbix",
					"version": "7.2"
				}
			},
			{
				"templateid": "10002",
				"host": "windows-template",
				"name": "Template OS Windows"
			}
		],
		"id": 1
	}`

	var response zabbix.TemplateGetResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	require.NoError(t, err)

	// Verify response structure
	require.Len(t, response.Result, 2)

	// Verify first template with all fields
	require.Equal(t, "10001", response.Result[0].TemplateID)
	require.Equal(t, "linux-template", response.Result[0].Host)
	require.Equal(t, "Template OS Linux", response.Result[0].Name)
	require.Equal(t, "Linux OS monitoring", response.Result[0].Description)
	require.Equal(t, "abc123", response.Result[0].UUID)
	require.Equal(t, "Zabbix", response.Result[0].Vendor.Name)
	require.Equal(t, "7.2", response.Result[0].Vendor.Version)

	// Verify second template with minimal fields
	require.Equal(t, "10002", response.Result[1].TemplateID)
	require.Equal(t, "windows-template", response.Result[1].Host)
	require.Equal(t, "Template OS Windows", response.Result[1].Name)
}

// TestTemplateGetJSONMarshaling tests that Template struct marshals correctly to JSON.
func TestTemplateGetJSONMarshaling(t *testing.T) {
	template := zabbix.Template{
		TemplateID:  "10001",
		Host:        "test-template",
		Name:        "Test Template",
		Description: "Test description",
		UUID:        "uuid-123",
	}
	template.Vendor.Name = "TestVendor"
	template.Vendor.Version = "1.0"

	jsonData, err := json.Marshal(template)
	require.NoError(t, err)

	// Verify JSON contains all fields with correct names
	jsonStr := string(jsonData)
	require.Contains(t, jsonStr, `"templateid":"10001"`)
	require.Contains(t, jsonStr, `"host":"test-template"`)
	require.Contains(t, jsonStr, `"name":"Test Template"`)
	require.Contains(t, jsonStr, `"description":"Test description"`)
	require.Contains(t, jsonStr, `"uuid":"uuid-123"`)
	require.Contains(t, jsonStr, `"vendor"`)
	require.Contains(t, jsonStr, `"name":"TestVendor"`)
	require.Contains(t, jsonStr, `"version":"1.0"`)
}

// TestTemplateGetIDFieldsAreStrings verifies all ID fields are strings, not integers.
func TestTemplateGetIDFieldsAreStrings(t *testing.T) {
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"templateid": "10001",
				"host": "template-host",
				"name": "Template Name"
			}
		],
		"id": 1
	}`

	var response zabbix.TemplateGetResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	require.NoError(t, err)

	// Verify templateid is a string, not an integer
	require.IsType(t, "", response.Result[0].TemplateID)
	require.Equal(t, "10001", response.Result[0].TemplateID)
}

// TestNewTemplateGetRequest tests the request constructor with options.
func TestNewTemplateGetRequest(t *testing.T) {
	req := zabbix.NewTemplateGetRequest(
		zabbix.WithTemplateGetTemplateIDs([]string{"10001", "10002"}),
		zabbix.WithTemplateGetOutput("extend"),
		zabbix.WithTemplateGetFilter(map[string]any{"name": []string{"Template OS Linux"}}),
	)

	require.Equal(t, "template.get", req.Method)
	require.Equal(t, "2.0", req.JSONRPC)
	require.Equal(t, []string{"10001", "10002"}, req.Params.TemplateIDs)
	require.Equal(t, "extend", req.Params.Output)
	require.NotNil(t, req.Params.Filter)
}

// TestTemplateGetWithFilters tests functional options for filter parameters.
func TestTemplateGetWithFilters(t *testing.T) {
	req := zabbix.NewTemplateGetRequest(
		zabbix.WithTemplateGetTemplateIDs([]string{"10001"}),
		zabbix.WithTemplateGetGroupIDs([]string{"5"}),
		zabbix.WithTemplateGetParentTemplateIDs([]string{"100"}),
		zabbix.WithTemplateGetHostIDs([]string{"20001"}),
		zabbix.WithTemplateGetGraphIDs([]string{"300"}),
		zabbix.WithTemplateGetItemIDs([]string{"400"}),
		zabbix.WithTemplateGetTriggerIDs([]string{"500"}),
	)

	require.Equal(t, []string{"10001"}, req.Params.TemplateIDs)
	require.Equal(t, []string{"5"}, req.Params.GroupIDs)
	require.Equal(t, []string{"100"}, req.Params.ParentTemplateIDs)
	require.Equal(t, []string{"20001"}, req.Params.HostIDs)
	require.Equal(t, []string{"300"}, req.Params.GraphIDs)
	require.Equal(t, []string{"400"}, req.Params.ItemIDs)
	require.Equal(t, []string{"500"}, req.Params.TriggerIDs)
}

// TestTemplateGetWithBooleanFlags tests boolean flag options.
func TestTemplateGetWithBooleanFlags(t *testing.T) {
	req := zabbix.NewTemplateGetRequest(
		zabbix.WithTemplateGetWithItems(true),
		zabbix.WithTemplateGetWithTriggers(true),
		zabbix.WithTemplateGetWithGraphs(true),
		zabbix.WithTemplateGetWithHTTPTests(true),
	)

	require.True(t, req.Params.WithItems)
	require.True(t, req.Params.WithTriggers)
	require.True(t, req.Params.WithGraphs)
	require.True(t, req.Params.WithHTTPTests)
}

// TestTemplateGetWithSelectOptions tests all select option functions.
func TestTemplateGetWithSelectOptions(t *testing.T) {
	req := zabbix.NewTemplateGetRequest(
		zabbix.WithTemplateGetSelectTags("extend"),
		zabbix.WithTemplateGetSelectHosts([]string{"hostid", "host"}),
		zabbix.WithTemplateGetSelectTemplateGroups("extend"),
		zabbix.WithTemplateGetSelectTemplates([]string{"templateid"}),
		zabbix.WithTemplateGetSelectParentTemplates("extend"),
		zabbix.WithTemplateGetSelectHTTPTests([]string{"httptestid"}),
		zabbix.WithTemplateGetSelectItems([]string{"itemid", "name"}),
		zabbix.WithTemplateGetSelectDiscoveries("extend"),
		zabbix.WithTemplateGetSelectTriggers([]string{"triggerid"}),
		zabbix.WithTemplateGetSelectGraphs([]string{"graphid"}),
		zabbix.WithTemplateGetSelectMacros("extend"),
		zabbix.WithTemplateGetSelectDashboards([]string{"dashboardid"}),
		zabbix.WithTemplateGetSelectValueMaps("extend"),
	)

	require.Equal(t, "extend", req.Params.SelectTags)
	require.Equal(t, []string{"hostid", "host"}, req.Params.SelectHosts)
	require.Equal(t, "extend", req.Params.SelectTemplateGroups)
	require.Equal(t, []string{"templateid"}, req.Params.SelectTemplates)
	require.Equal(t, "extend", req.Params.SelectParentTemplates)
	require.Equal(t, []string{"httptestid"}, req.Params.SelectHTTPTests)
	require.Equal(t, []string{"itemid", "name"}, req.Params.SelectItems)
	require.Equal(t, "extend", req.Params.SelectDiscoveries)
	require.Equal(t, []string{"triggerid"}, req.Params.SelectTriggers)
	require.Equal(t, []string{"graphid"}, req.Params.SelectGraphs)
	require.Equal(t, "extend", req.Params.SelectMacros)
	require.Equal(t, []string{"dashboardid"}, req.Params.SelectDashboards)
	require.Equal(t, "extend", req.Params.SelectValueMaps)
}

// TestNewGetAllTemplatesRequest tests the convenience constructor.
func TestNewGetAllTemplatesRequest(t *testing.T) {
	req := zabbix.NewGetAllTemplatesRequest()

	require.Equal(t, "template.get", req.Method)
	require.Equal(t, "extend", req.Params.Output)
}

// TestTemplateVendorStructUnmarshaling tests unmarshaling of nested vendor struct.
func TestTemplateVendorStructUnmarshaling(t *testing.T) {
	jsonResponse := `{
		"jsonrpc": "2.0",
		"result": [
			{
				"templateid": "10001",
				"host": "vendor-template",
				"name": "Template with Vendor",
				"vendor": {
					"name": "Acme Corp",
					"version": "2.5.1"
				}
			}
		],
		"id": 1
	}`

	var response zabbix.TemplateGetResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	require.NoError(t, err)

	require.Len(t, response.Result, 1)
	require.Equal(t, "Acme Corp", response.Result[0].Vendor.Name)
	require.Equal(t, "2.5.1", response.Result[0].Vendor.Version)
}

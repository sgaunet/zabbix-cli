package zabbix_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	t.Parallel()

	t.Run("Import success", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request
				var req map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&req)
				if err != nil {
					t.Fatalf("Failed to decode request: %v", err)
				}

				// Verify the request
				require.Equal(t, "2.0", req["jsonrpc"])
				require.Equal(t, "configuration.import", req["method"])
				require.Equal(t, "auth_token", req["auth"])

				params, ok := req["params"].(map[string]interface{})
				require.True(t, ok, "params should be a map")
				require.Equal(t, "yaml", params["format"])
				require.Equal(t, "test_import_data", params["source"])

				// Return success response
				res := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"result":  true,
					"id":      1,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(res)
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import
		success, err := z.Import(context.Background(), "test_import_data")
		require.NoError(t, err)
		require.True(t, success)
	})

	t.Run("Import HTTP error", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request, return HTTP error
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"error":   map[string]interface{}{"code": -32603, "message": "Internal error"},
					"id":      1,
				})
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import - should fail with HTTP error
		_, err = z.Import(context.Background(), "test_import_data")
		require.Error(t, err)
		require.Contains(t, err.Error(), "status code not OK: 500")
	})

	t.Run("Import request error", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request, close the connection
				conn, _, _ := w.(http.Hijacker).Hijack()
				conn.Close()
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import - should fail with connection error
		_, err = z.Import(context.Background(), "test_import_data")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot do request")
	})

	t.Run("Import API error", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request, return API error
				w.Header().Set("Content-Type", "application/json")
				errResp := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"error": map[string]interface{}{
						"code":    -32602,
						"message": "Invalid params.",
						"data":    "Incorrect parameter \"source\"",
					},
					"id": 1,
				}
				json.NewEncoder(w).Encode(errResp)
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import - should fail with API error
		_, err = z.Import(context.Background(), "test_import_data")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Zabbix API error")
		require.Contains(t, err.Error(), "Invalid params")
		require.Contains(t, err.Error(), "Incorrect parameter \"source\"")
	})

	t.Run("Import response returns boolean correctly", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request - return boolean true
				res := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"result":  true,
					"id":      1,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(res)
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import
		success, err := z.Import(context.Background(), "valid_yaml_data")
		require.NoError(t, err)
		require.True(t, success)
		require.IsType(t, true, success) // Verify it's a boolean
	})

	t.Run("Import malformed JSON response", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request, return malformed JSON
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"invalid json`))
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import - should fail with unmarshal error
		_, err = z.Import(context.Background(), "test_data")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal response")
	})

	t.Run("Import with empty source", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request with empty source, return API error
				errResp := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"error": map[string]interface{}{
						"code":    -32602,
						"message": "Invalid params.",
						"data":    "Import data is missing",
					},
					"id": 1,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(errResp)
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import with empty source - should fail
		_, err = z.Import(context.Background(), "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Zabbix API error")
	})

	t.Run("Import with invalid YAML format", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and import requests
		var loginCalled bool
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For login request
				if !loginCalled {
					loginCalled = true
					res := zabbix.LoginResponse{
						JSONRPC: zabbix.JSONRPC,
						Result:  "auth_token",
						ID:      1,
					}
					resJSON, _ := json.Marshal(res)
					w.Header().Add("Content-Type", "application/json")
					w.Write(resJSON)
					return
				}

				// For import request with invalid YAML, return API error
				errResp := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"error": map[string]interface{}{
						"code":    -32602,
						"message": "Invalid params.",
						"data":    "Invalid import format",
					},
					"id": 1,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(errResp)
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Test import with invalid YAML - should fail
		_, err = z.Import(context.Background(), "invalid: yaml: format:")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid import format")
	})
}

func TestImportRulesMarshaling(t *testing.T) {
	t.Parallel()

	t.Run("Rules with all options true", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test_source")
		payload, err := json.Marshal(req)
		require.NoError(t, err)
		require.NotEmpty(t, payload)

		// Verify all rule types are present
		require.Contains(t, string(payload), "discoveryRules")
		require.Contains(t, string(payload), "graphs")
		require.Contains(t, string(payload), "host_groups")
		require.Contains(t, string(payload), "template_groups")
		require.Contains(t, string(payload), "hosts")
		require.Contains(t, string(payload), "httptests")
		require.Contains(t, string(payload), "images")
		require.Contains(t, string(payload), "items")
		require.Contains(t, string(payload), "maps")
		require.Contains(t, string(payload), "mediaTypes")
		require.Contains(t, string(payload), "templateLinkage")
		require.Contains(t, string(payload), "templates")
		require.Contains(t, string(payload), "templateDashboards")
		require.Contains(t, string(payload), "triggers")
		require.Contains(t, string(payload), "valueMaps")

		// Verify format and source
		require.Contains(t, string(payload), `"format":"yaml"`)
		require.Contains(t, string(payload), `"source":"test_source"`)
	})

	t.Run("Host groups rules have no deleteMissing", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		// Parse the payload to check host_groups structure
		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		hostGroups := rules["host_groups"].(map[string]interface{})

		// host_groups should only have createMissing and updateExisting
		require.Contains(t, hostGroups, "createMissing")
		require.Contains(t, hostGroups, "updateExisting")
		require.NotContains(t, hostGroups, "deleteMissing")
	})

	t.Run("Template groups rules have no deleteMissing", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		templateGroups := rules["template_groups"].(map[string]interface{})

		// template_groups should only have createMissing and updateExisting
		require.Contains(t, templateGroups, "createMissing")
		require.Contains(t, templateGroups, "updateExisting")
		require.NotContains(t, templateGroups, "deleteMissing")
	})

	t.Run("TemplateLinkage has correct fields", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		templateLinkage := rules["templateLinkage"].(map[string]interface{})

		// templateLinkage should have createMissing and deleteMissing (not updateExisting)
		require.Contains(t, templateLinkage, "createMissing")
		require.Contains(t, templateLinkage, "deleteMissing")
		require.NotContains(t, templateLinkage, "updateExisting")

		// Verify the JSON uses correct field name "deleteMissing" not "deleteExisting"
		require.Contains(t, string(payload), `"deleteMissing"`)
	})

	t.Run("Discovery rules have all three boolean fields", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		discoveryRules := rules["discoveryRules"].(map[string]interface{})

		// discoveryRules should have all three fields
		require.Contains(t, discoveryRules, "createMissing")
		require.Contains(t, discoveryRules, "updateExisting")
		require.Contains(t, discoveryRules, "deleteMissing")
	})

	t.Run("Items have all three boolean fields", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		items := rules["items"].(map[string]interface{})

		// items should have all three fields
		require.Contains(t, items, "createMissing")
		require.Contains(t, items, "updateExisting")
		require.Contains(t, items, "deleteMissing")
	})

	t.Run("ValueMaps have all three boolean fields", func(t *testing.T) {
		t.Parallel()
		req := zabbix.NewConfigurationImportRequest("test")
		payload, err := json.Marshal(req)
		require.NoError(t, err)

		var parsed map[string]interface{}
		err = json.Unmarshal(payload, &parsed)
		require.NoError(t, err)

		params := parsed["params"].(map[string]interface{})
		rules := params["rules"].(map[string]interface{})
		valueMaps := rules["valueMaps"].(map[string]interface{})

		// valueMaps should have all three fields
		require.Contains(t, valueMaps, "createMissing")
		require.Contains(t, valueMaps, "updateExisting")
		require.Contains(t, valueMaps, "deleteMissing")
	})
}

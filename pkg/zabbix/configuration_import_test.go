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
}

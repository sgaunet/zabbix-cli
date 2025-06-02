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

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Login success", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				res := zabbix.LoginResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "auth_token",
					ID:      1,
				}
				resJSON, err := json.Marshal(res)
				if err != nil {
					t.Fatal(err)
				}
				w.Header().Add("Content-Type", "application/json")
				fmt.Fprintln(w, string(resJSON))
			}))
		defer ts.Close()

		// get request
		client := ts.Client()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)
		require.Equal(t, ts.URL, z.APIEndpoint)
		require.Equal(t, "user", z.User)
		require.Equal(t, "password", z.Password)

		z.SetHTTPClient(client)
		err := z.Login(context.Background())
		require.NoError(t, err)
		require.Equal(t, "auth_token", z.Auth())
	})

	t.Run("Login failure", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				res := zabbix.LoginResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "",
					ID:      1,
				}
				resJSON, err := json.Marshal(res)
				if err != nil {
					t.Fatal(err)
				}
				w.Header().Add("Content-Type", "application/json")
				fmt.Fprintln(w, string(resJSON))
			}))
		defer ts.Close()

		// get request
		client := ts.Client()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)

		z.SetHTTPClient(client)
		err := z.Login(context.Background())
		require.Error(t, err)
		require.Empty(t, z.Auth())
	})

	t.Run("Login request error", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
		defer ts.Close()

		// get request
		client := ts.Client()

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)

		z.SetHTTPClient(client)
		err := z.Login(context.Background())
		require.Error(t, err)
		require.Empty(t, z.Auth())
	})

	t.Run("Login request to unavailable endpoint", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
		client := ts.Client()
		ts.Close() // close the server

		z := zabbix.New("user", "password", ts.URL)
		require.NotNil(t, z)

		z.SetHTTPClient(client)
		err := z.Login(context.Background())
		require.Error(t, err)
		require.Empty(t, z.Auth())
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()

	t.Run("Logout success", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and content type
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Parse request body
				var req map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&req)
				require.NoError(t, err)

				// Verify request fields
				require.Equal(t, zabbix.JSONRPC, req["jsonrpc"])
				require.Equal(t, "user.logout", req["method"])
				require.Equal(t, "auth_token", req["auth"])

				// Send success response
				res := map[string]interface{}{
					"jsonrpc": zabbix.JSONRPC,
					"result":  true,
					"id":      req["id"],
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(res)
			}),
		)
		defer ts.Close()

		// Create a test server that will handle the login request
		loginTS := httptest.NewTLSServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				res := zabbix.LoginResponse{
					JSONRPC: zabbix.JSONRPC,
					Result:  "auth_token",
					ID:      1,
				}
				resJSON, _ := json.Marshal(res)
				w.Header().Add("Content-Type", "application/json")
				w.Write(resJSON)
			}),
		)
		defer loginTS.Close()

		client := loginTS.Client()
		z := zabbix.New("user", "password", loginTS.URL)
		z.SetHTTPClient(client)

		// Login to set the auth token
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Update the client to use the logout test server
		z.SetHTTPClient(ts.Client())

		err = z.Logout(context.Background())
		require.NoError(t, err)
	})

	t.Run("Logout HTTP error", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and logout requests
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

				// For logout request
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

		// Login should succeed
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Logout should fail with HTTP error
		err = z.Logout(context.Background())
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected status code: 500")
	})

	t.Run("Logout request error", func(t *testing.T) {
		t.Parallel()
		// Create a test server that will handle both login and logout requests
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

				// For logout request, close the connection
				conn, _, _ := w.(http.Hijacker).Hijack()
				conn.Close()
			}),
		)
		defer ts.Close()

		// Create client with the test server URL
		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		// Login should succeed
		err := z.Login(context.Background())
		require.NoError(t, err)

		// Logout should fail with connection error
		err = z.Logout(context.Background())
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot do request")
	})
}

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
		require.EqualError(t, err, "unexpected status code 500: wrong http code")
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

func TestHostGroupGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("HostGroupGet success", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			var req zabbix.HostGroupGetRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			require.Equal(t, "test_auth_token", req.Auth)
			require.Equal(t, 123, req.ID)

			resp := zabbix.HostGroupGetResponse{
				JSONRPC: zabbix.JSONRPC,
				Result: []zabbix.HostGroup{
					{GroupID: "1", Name: "Linux servers", Internal: "0"},
					{GroupID: "2", Name: "Windows servers", Internal: "0"},
				},
				ID: 123,
			}
			jsonResp, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(jsonResp))
		}))
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth("test_auth_token"),
			zabbix.WithHostGroupGetID(123),
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Len(t, response.Result, 2)
		require.Equal(t, "Linux servers", response.Result[0].Name)
		require.Nil(t, response.Error)
	})

	t.Run("HostGroupGet API error", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := zabbix.HostGroupGetResponse{
				JSONRPC: zabbix.JSONRPC,
				Error: &zabbix.Error{
					Code:    -32602,
					Message: "Invalid params.",
					Data:    "Invalid parameter \"/output\": value must be one of \"extend\", \"refer\", \"selectParent\", ...",
				},
				ID: 123,
			}
			jsonResp, _ := json.Marshal(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // Zabbix API errors are often returned with HTTP 200 OK
			fmt.Fprintln(w, string(jsonResp))
		}))
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth("test_auth_token"),
			zabbix.WithHostGroupGetID(123),
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		require.Error(t, err)
		require.Nil(t, response)
		zbxErr, ok := err.(*zabbix.Error)
		require.True(t, ok, "error should be of type *zabbix.Error")
		require.Equal(t, -32602, zbxErr.Code)
		require.Contains(t, zbxErr.Data, "Invalid parameter")
	})

	t.Run("HostGroupGet HTTP error", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"jsonrpc": "%s", "error": {"code": -32099, "message": "Internal error", "data": "Server is experiencing issues."}, "id": 123}`, zabbix.JSONRPC)
		}))
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth("test_auth_token"),
			zabbix.WithHostGroupGetID(123),
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		require.Error(t, err)
		require.Nil(t, response)
		require.EqualError(t, err, "Zabbix API error (Code: -32099): Internal error - Server is experiencing issues.")
		// Check if the Zabbix error is parsed from the body even on HTTP error
		require.Contains(t, err.Error(), "Internal error") 
	})

	t.Run("HostGroupGet malformed JSON response", func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "this is not valid json")
		}))
		defer ts.Close()

		z := zabbix.New("user", "password", ts.URL)
		z.SetHTTPClient(ts.Client())

		hgRequest := zabbix.NewGetAllHostGroupsRequest(
			zabbix.WithHostGroupGetAuth("test_auth_token"),
			zabbix.WithHostGroupGetID(123),
		)

		response, err := z.HostGroupGet(ctx, hgRequest)
		require.Error(t, err)
		require.Nil(t, response)
		require.Contains(t, err.Error(), "cannot unmarshal hostgroup.get response")
	})
}


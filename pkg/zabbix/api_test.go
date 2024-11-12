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
				res := zabbix.ZbxLoginResponse{
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
				res := zabbix.ZbxLoginResponse{
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

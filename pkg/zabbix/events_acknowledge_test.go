package zabbix_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	zabbix "github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

func TestAcknowledgeEvents(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result: struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101, 102}},
				ID:      1,
			}
			resJSON, err := json.Marshal(res)
			require.NoError(t, err)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())
		// Normally, zabbix.Client initializes auth/id during Login, but here we are testing AcknowledgeEvents directly.
		// If AcknowledgeEvents requires authentication, you must call Login first or mock it. For this test, we assume the test server will accept any auth.
		ids, err := z.AcknowledgeEvents(context.Background(), []string{"1", "2"})
		require.NoError(t, err)
		require.Equal(t, []int{101, 102}, ids)
	})

	t.Run("network error", func(t *testing.T) {
		z := zabbix.New("user", "pass", "http://127.0.0.1:0") // invalid port to force error
		_, err := z.AcknowledgeEvents(context.Background(), []string{"1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot do request")
	})

	t.Run("non-200 status", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"error": "bad request"}`)
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())
		_, err := z.AcknowledgeEvents(context.Background(), []string{"1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "status code not OK")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "not json")
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())
		_, err := z.AcknowledgeEvents(context.Background(), []string{"1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal response")
	})

	t.Run("zabbix error message", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				ID:      1,
				ErrorMsg: zabbix.ErrorMsg{Code: 123, Message: "fail", Data: "details"},
			}
			resJSON, err := json.Marshal(res)
			require.NoError(t, err)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()
		
		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())
		_, err := z.AcknowledgeEvents(context.Background(), []string{"1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "error message")
	})
}

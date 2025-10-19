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
			res := struct {
				JSONRPC string `json:"jsonrpc"`
				ID      int    `json:"id"`
				Error   struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Data    string `json:"data"`
				} `json:"error"`
			}{
				JSONRPC: zabbix.JSONRPC,
				ID:      1,
				Error: struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Data    string `json:"data"`
				}{Code: 123, Message: "fail", Data: "details"},
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
		require.Contains(t, err.Error(), "Zabbix API error")
	})
}

func TestAcknowledgeEventsWithOptions(t *testing.T) {
	t.Parallel()

	t.Run("with message option", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithMessage("Test message"),
		)
		require.NoError(t, err)
		require.Equal(t, "Test message", receivedParams.Message)
	})

	t.Run("with severity option", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithSeverity(zabbix.High),
		)
		require.NoError(t, err)
		require.Equal(t, zabbix.High, receivedParams.Severity)
	})

	t.Run("with actions - acknowledge and add message", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithActions(zabbix.Acknowledge, zabbix.AddMessage),
			zabbix.WithMessage("Problem acknowledged"),
		)
		require.NoError(t, err)
		require.Equal(t, 6, receivedParams.Action) // 2 + 4
		require.Equal(t, "Problem acknowledged", receivedParams.Message)
	})

	t.Run("with actions - change severity", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithActions(zabbix.ChangeSeverity),
			zabbix.WithSeverity(zabbix.Disaster),
		)
		require.NoError(t, err)
		require.Equal(t, 8, receivedParams.Action)
		require.Equal(t, zabbix.Disaster, receivedParams.Severity)
	})

	t.Run("with multiple actions and options", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1", "2"},
			zabbix.WithActions(zabbix.CloseProblem, zabbix.Acknowledge, zabbix.AddMessage, zabbix.ChangeSeverity),
			zabbix.WithMessage("Resolving with high severity"),
			zabbix.WithSeverity(zabbix.High),
		)
		require.NoError(t, err)
		require.Equal(t, 15, receivedParams.Action) // 1 + 2 + 4 + 8
		require.Equal(t, "Resolving with high severity", receivedParams.Message)
		require.Equal(t, zabbix.High, receivedParams.Severity)
		require.Equal(t, []string{"1", "2"}, receivedParams.Eventids)
	})

	t.Run("with suppress action", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithActions(zabbix.Suppress),
		)
		require.NoError(t, err)
		require.Equal(t, 32, receivedParams.Action)
	})

	t.Run("with unsuppress action", func(t *testing.T) {
		var receivedParams zabbix.EventsAcknowledgeParams
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Params zabbix.EventsAcknowledgeParams `json:"params"`
			}
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			receivedParams = req.Params

			res := zabbix.EventAcknowledgeResponse{
				JSONRPC: zabbix.JSONRPC,
				Result:  struct{ Eventids []int `json:"eventids"` }{Eventids: []int{101}},
				ID:      1,
			}
			resJSON, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(w, string(resJSON))
		}))
		defer ts.Close()

		z := zabbix.New("user", "pass", ts.URL)
		z.SetHTTPClient(ts.Client())

		_, err := z.AcknowledgeEvents(
			context.Background(),
			[]string{"1"},
			zabbix.WithActions(zabbix.Unsuppress),
		)
		require.NoError(t, err)
		require.Equal(t, 64, receivedParams.Action)
	})
}

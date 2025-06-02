package zabbix_test

import (
	"context"
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
			res := `{
				"jsonrpc": "2.0",
				"result": [
					{
						"templateid": "10001",
						"name": "Template OS Linux"
					}
				],
				"id": 1
			}`
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
		Result: []struct {
			TemplateID string `json:"templateID"`
			Name       string `json:"name"`
		}{
			{TemplateID: "10001", Name: "Template OS Linux"},
			{TemplateID: "10002", Name: "Template OS Windows"},
		},
	}
	ids := resp.GetTemplateID()
	require.Equal(t, []string{"10001", "10002"}, ids)
}

func TestGetTemplateName_Internal(t *testing.T) {
	resp := &zabbix.TemplateGetResponse{
		Result: []struct {
			TemplateID string `json:"templateID"`
			Name       string `json:"name"`
		}{
			{TemplateID: "10001", Name: "Template OS Linux"},
			{TemplateID: "10002", Name: "Template OS Windows"},
		},
	}
	names := resp.GetTemplateName()
	require.Equal(t, []string{"Template OS Linux", "Template OS Windows"}, names)
}

package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// defaultTimeout is the default timeout for the HTTP client
const defaultTimeout = 5 * time.Second

const methodUserLogin = "user.login"
const methodUserLogout = "user.logout"

// New creates a new ZabbixAPI object
// and logs in to the Zabbix API
// It returns the ZabbixAPI object
// and an error if any
// The default timeout is 5 seconds
// Don't forget to call Logout() to logout
func New(user, password, apiEndpoint string) (ZabbixAPI, error) {
	z := ZabbixAPI{
		APIEndpoint: apiEndpoint,
		User:        user,
		Password:    password,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	data := zbxRequestLogin{
		JSONRPC: JSONRPC,
		Method:  methodUserLogin,
		Params: zbxParams{
			UserName: z.User,
			Password: z.Password,
		},
		ID: generateUniqueID(),
	}

	statusCode, resp, err := z.postRequest(context.Background(), data)
	if err != nil {
		return z, fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return z, fmt.Errorf("unexpected status code: %d (%w)", statusCode, ErrWrongHTTPCode)
	}
	var zbxResp zbxLoginResponse
	err = json.Unmarshal(resp, &zbxResp)
	if err != nil {
		return z, fmt.Errorf("cannot unmarshal response: %w - %s", err, string(resp))
	}
	if zbxResp.Result == "" {
		return z, fmt.Errorf("cannot login: %w", ErrEmptyResult)
	}
	z.auth = zbxResp.Result
	return z, nil
}

// Logout logs out from the Zabbix API
func (z *ZabbixAPI) Logout(ctx context.Context) error {
	data := zbxRequestLogout{
		JSONRPC: JSONRPC,
		Method:  methodUserLogout,
		Params:  make(map[string]string),
		ID:      z.id,
		Auth:    z.auth,
	}

	statusCode, _, err := z.postRequest(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d (%w)", statusCode, ErrWrongHTTPCode)
	}
	return nil
}

// Auth returns the auth token
// that is used to authenticate
// This token is initialized during the login process
func (z *ZabbixAPI) Auth() string {
	return z.auth
}

// postRequest sends a POST request to the Zabbix API
// It returns the status code, the response body and an error if any
func (z *ZabbixAPI) postRequest(ctx context.Context, payload interface{}) (int, []byte, error) {
	return z.request(ctx, http.MethodPost, payload)
}

// func (z *ZabbixAPI) getRequest(ctx context.Context, payload interface{}) (int, []byte, error) {
// 	return z.request(ctx, http.MethodGet, c)
// }

// request sends a request to the Zabbix API
// It returns the status code, the response body and an error if any
func (z *ZabbixAPI) request(ctx context.Context, method string, payload interface{}) (int, []byte, error) {
	postBody, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(method, z.APIEndpoint, responseBody)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot read response body: %w", err)
	}
	return resp.StatusCode, body, nil
}

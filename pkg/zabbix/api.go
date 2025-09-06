// Package zabbix provides the core Zabbix API client implementation.
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

// defaultTimeout is the default timeout for the HTTP client.
const defaultTimeout = 5 * time.Second

const methodUserLogin = "user.login"
const methodUserLogout = "user.logout"

// New creates a new Client object
// The default timeout is 5 seconds.
func New(user, password, apiEndpoint string) Client {
	return Client{
		APIEndpoint: apiEndpoint,
		User:        user,
		Password:    password,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SetHTTPClient sets the HTTP client.
func (z *Client) SetHTTPClient(client *http.Client) {
	z.client = client
}

// Login logs in to the Zabbix API.
// Don't forget to call Logout() to logout.
func (z *Client) Login(ctx context.Context) error {
	data := LoginRequest{
		JSONRPC: JSONRPC,
		Method:  methodUserLogin,
		Params: Params{
			UserName: z.User,
			Password: z.Password,
		},
		ID: generateUniqueID(),
	}

	statusCode, resp, err := z.postRequest(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %w", statusCode, ErrWrongHTTPCode)
	}
	var zbxResp LoginResponse
	err = json.Unmarshal(resp, &zbxResp)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response: %w - %s", err, string(resp))
	}
	if zbxResp.Result == "" {
		return fmt.Errorf("login failed, empty auth token: %w", ErrEmptyResult)
	}
	z.auth = zbxResp.Result
	return nil
}

// Logout logs out from the Zabbix API.
func (z *Client) Logout(ctx context.Context) error {
	data := LogoutRequest{
		JSONRPC: JSONRPC,
		Method:  methodUserLogout,
		Params:  make([]interface{}, 0),
		ID:      generateUniqueID(),
		Auth:    z.auth,
	}

	statusCode, _, err := z.postRequest(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot do request: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %w", statusCode, ErrWrongHTTPCode)
	}
	return nil
}

// HostGroupGet sends a hostgroup.get request to the Zabbix API.
// The request object should be fully populated by the caller, including Auth and ID.
func (z *Client) HostGroupGet(ctx context.Context, request *HostGroupGetRequest) (*HostGroupGetResponse, error) {
	statusCode, respBody, err := z.postRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API request failed for hostgroup.get: %w", err)
	}

	var response HostGroupGetResponse
	if err := handleRawResponse(statusCode, respBody, "hostgroup.get", &response); err != nil {
		return nil, err
	}

	if response.Error != nil && response.Error.Code != 0 {
		return nil, response.Error
	}

	return &response, nil
}

// Auth returns the auth token
// that is used to authenticate.
// This token is initialized during the login process.
func (z *Client) Auth() string {
	return z.auth
}

// postRequest sends a POST request to the Zabbix API.
// It returns the status code, the response body and an error if any.
func (z *Client) postRequest(ctx context.Context, payload interface{}) (int, []byte, error) {
	return z.request(ctx, http.MethodPost, payload)
}

// request sends a request to the Zabbix API.
// It returns the status code, the response body and an error if any.
func (z *Client) request(ctx context.Context, method string, payload interface{}) (int, []byte, error) {
	postBody, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequestWithContext(ctx, method, z.APIEndpoint, responseBody)
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

package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// defaultTimeout is the default timeout for the HTTP client
const defaultTimeout = 5 * time.Second

const MethodUserLogin = "user.login"
const MethodUserLogout = "user.logout"

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
		Method:  MethodUserLogin,
		Params: zbxParams{
			UserName: z.User,
			Password: z.Password,
		},
		ID: 1,
	}

	postBody, err := json.Marshal(data)
	if err != nil {
		return z, fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, z.APIEndpoint, responseBody)
	if err != nil {
		return z, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(context.TODO())
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return z, fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return z, err
	}
	var zbxResp zbxLoginResponse
	err = json.Unmarshal(body, &zbxResp)
	if err != nil {
		return z, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	// fmt.Println(string(body))
	if zbxResp.Result == "" {
		return z, fmt.Errorf("cannot login: %w", errors.New("empty result"))
	}
	z.auth = zbxResp.Result
	return z, err
}

// Logout logs out from the Zabbix API
func (z *ZabbixAPI) Logout(ctx context.Context) error {
	data := zbxRequestLogout{
		JSONRPC: JSONRPC,
		Method:  MethodUserLogout,
		Params:  make(map[string]string),
		ID:      1,
		Auth:    z.auth,
	}

	postBody, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal data: %w", err)
	}
	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", z.APIEndpoint, responseBody)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// Auth returns the auth token
// that is used to authenticate
// This token is initialized during the login process
func (z *ZabbixAPI) Auth() string {
	return z.auth
}

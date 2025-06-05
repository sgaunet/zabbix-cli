package zabbix

import "fmt"

// APIError represents a structured error from the Zabbix API
// It implements the error interface
var (
	// ErrUnexpectedResponse is returned when the API returns an unexpected response format
	ErrUnexpectedResponse = NewAPIError("unexpected API response format")
)

// APIError represents an error from the Zabbix API
// It implements the error interface
type APIError struct {
	msg string
}

// NewAPIError creates a new APIError with the given message
func NewAPIError(msg string) *APIError {
	return &APIError{msg: msg}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.msg
}

// WithStatus adds HTTP status code to the error
func (e *APIError) WithStatus(statusCode int) *APIError {
	e.msg = fmt.Sprintf("%s (status: %d)", e.msg, statusCode)
	return e
}

// WithBody adds response body to the error
func (e *APIError) WithBody(body []byte) *APIError {
	e.msg = fmt.Sprintf("%s: %s", e.msg, string(body))
	return e
}

// WithMethod adds API method name to the error
func (e *APIError) WithMethod(method string) *APIError {
	e.msg = fmt.Sprintf("API request for %s %s", method, e.msg)
	return e
}

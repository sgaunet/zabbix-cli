package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleRawResponse processes the raw HTTP response from a Zabbix API call.
// It handles non-OK HTTP status codes by attempting to parse a Zabbix error object.
// If the status is OK, it unmarshals the response body into responsePtr.
// methodName is used for enriching error messages.
func handleRawResponse(statusCode int, respBody []byte, methodName string, responsePtr interface{}) error {
	if statusCode != http.StatusOK {
		var zbxErrorResponse struct {
			Error *Error `json:"error,omitempty"` // Assumes Error type is defined in common.go
		}
		// Attempt to parse a Zabbix-specific error from the body
		if unmarshalErr := json.Unmarshal(respBody, &zbxErrorResponse); unmarshalErr == nil && zbxErrorResponse.Error != nil && zbxErrorResponse.Error.Code != 0 {
			return zbxErrorResponse.Error // Return the Zabbix-specific error
		}
		// If not a Zabbix-specific error or parsing failed, return a generic error using ErrUnexpectedResponse from errors.go
		return ErrUnexpectedResponse.
			WithMethod(methodName).
			WithStatus(statusCode).
			WithBody(respBody)
	}

	// Status is OK, unmarshal the response into the provided pointer
	err := json.Unmarshal(respBody, responsePtr)
	if err != nil {
		return fmt.Errorf("cannot unmarshal %s response: %w - %s", methodName, err, string(respBody))
	}
	return nil
}

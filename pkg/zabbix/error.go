package zabbix

import "errors"

// ErrEmptyResult is returned when the result is empty.
var ErrEmptyResult = errors.New("empty result")

// ErrWrongHTTPCode is returned when the HTTP code is wrong.
var ErrWrongHTTPCode = errors.New("wrong http code")

// ErrorMsg is the error message returned by the Zabbix API.
type ErrorMsg struct { //nolint:errname
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Error returns the error message.
func (e *ErrorMsg) Error() string {
	return e.Message + " " + e.Data
}

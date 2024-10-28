package zabbix

import "errors"

var ErrEmptyResult = errors.New("empty result")

// ErrorMsg is the error message returned by the Zabbix API
type ErrorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *ErrorMsg) Error() string {
	return e.Message + " " + e.Data
}

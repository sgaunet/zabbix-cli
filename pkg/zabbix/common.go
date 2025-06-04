package zabbix

import (
	"fmt"
)

// Error represents an error object from the Zabbix API.
// It is used to deserialize error responses from Zabbix.
type Error struct {
	Code    int    `json:"code"`    // Application-level error code.
	Message string `json:"message"` // Error message.
	Data    string `json:"data"`    // Detailed error information.
}

// Error implements the error interface for the Error struct.
func (e Error) Error() string {
	return fmt.Sprintf("Zabbix API error (Code: %d): %s - %s", e.Code, e.Message, e.Data)
}

// CommonGetParams represents common parameters for Zabbix API 'get' methods.
// This can be embedded in other specific GetParams structs.
type CommonGetParams struct {
	Output             interface{}            `json:"output,omitempty"`             // "extend", "select", etc.
	Limit              int                    `json:"limit,omitempty"`              // Max number of results to return.
	Filter             map[string]interface{} `json:"filter,omitempty"`             // Key-value pairs for filtering.
	SortField          []string               `json:"sortfield,omitempty"`          // Fields to sort by.
	SortOrder          []string               `json:"sortorder,omitempty"`          // Sort order ("ASC" or "DESC").
	Search             map[string]interface{} `json:"search,omitempty"`             // Key-value pairs for searching.
	SearchByAny        bool                   `json:"searchByAny,omitempty"`        // If true, search by any field.
	SearchWildcardsEnabled bool                `json:"searchWildcardsEnabled,omitempty"` // If true, enable wildcards in search.
	StartSearch        bool                   `json:"startSearch,omitempty"`        // If true, search from the beginning of the string.
	ExcludeSearch      bool                   `json:"excludeSearch,omitempty"`      // If true, exclude results matching the search.
	PreserveKeys       bool                   `json:"preservekeys,omitempty"`       // If true, preserve keys in the output.
	CountOutput        bool                   `json:"countOutput,omitempty"`        // If true, return the count of results.
	Editable           bool                   `json:"editable,omitempty"`           // If true, return only objects that are available for editing.
}

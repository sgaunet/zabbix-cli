package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const MethodEventAcknowledge = "event.acknowledge"

type zbxEventsAcknowledgeParams struct {
	Eventids []string `json:"eventids"`

	// Event update action(s).
	// Possible bitmap values are:
	// 1 - close problem;
	// 2 - acknowledge event;
	// 4 - add message;
	// 8 - change severity;
	// 16 - unacknowledge event.

	// This is a bitmask field; any sum of possible bitmap values is acceptable (for example, 6 for acknowledge event and add message).
	Action  int    `json:"action"`
	Message string `json:"message"`

	// 	New severity for events.
	Severity Severity `json:"severity"` // Required, if action contains 'change severity' flag.
}

type zbxEventAcknowledge struct {
	JSONRPC string                     `json:"jsonrpc"`
	Method  string                     `json:"method"`
	Params  zbxEventsAcknowledgeParams `json:"params"`
	Auth    string                     `json:"auth"`
	ID      int                        `json:"id"`
}

type zbxEventAcknowledgeOption func(*zbxEventAcknowledge)

// WithSeverity sets the severity of the event.
func WithSeverity(severity Severity) zbxEventAcknowledgeOption {
	return func(e *zbxEventAcknowledge) {
		e.Params.Severity = severity
	}
}

// WithAction sets the action to perform on the event.
func WithActions(action ...EventAction) zbxEventAcknowledgeOption {
	actions := NewEventAction(action...)
	return func(e *zbxEventAcknowledge) {
		e.Params.Action = actions
	}
}

// WithMessage sets the message to add to the event.
func WithMessage(message string) zbxEventAcknowledgeOption {
	return func(e *zbxEventAcknowledge) {
		e.Params.Message = message
	}
}

// newEventAcknowledgeRequest creates a new event acknowledge request.
func newEventAcknowledgeRequest(eventids []string, opts ...zbxEventAcknowledgeOption) *zbxEventAcknowledge {
	req := &zbxEventAcknowledge{
		JSONRPC: JSONRPC,
		Method:  MethodEventAcknowledge,
		Params: zbxEventsAcknowledgeParams{
			Eventids: eventids,
		},
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

type EventAcknowledgeResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Eventids []int `json:"eventids"`
	} `json:"result"`
	ID       int      `json:"id"`
	ErrorMsg ErrorMsg `json:"error,omitempty"`
}

func (z *ZabbixAPI) AcknowledgeEvents(ctx context.Context, eventsID []string, opts ...zbxEventAcknowledgeOption) ([]int, error) {
	payload := newEventAcknowledgeRequest(eventsID, opts...)
	payload.Auth = z.auth
	payload.ID = z.id
	statusCode, body, err := z.postRequest(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot do request: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not OK: %d - %s (%w)", statusCode, string(body), ErrWrongHTTPCode)
	}
	var res EventAcknowledgeResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	if res.ErrorMsg != (ErrorMsg{}) {
		return nil, fmt.Errorf("error message: %w", &res.ErrorMsg)
	}
	return res.Result.Eventids, nil
}

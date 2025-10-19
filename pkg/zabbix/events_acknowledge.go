package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MethodEventAcknowledge is the Zabbix API method for event acknowledgement.
const MethodEventAcknowledge = "event.acknowledge"

// EventsAcknowledgeParams contains parameters for acknowledging Zabbix events.
type EventsAcknowledgeParams struct {
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

// EventAcknowledgeRequest represents a request to acknowledge Zabbix events.
type EventAcknowledgeRequest struct {
	JSONRPC string                  `json:"jsonrpc"`
	Method  string                  `json:"method"`
	Params  EventsAcknowledgeParams `json:"params"`
	Auth    string                  `json:"auth"`
	ID      int                     `json:"id"`
}

// EventAcknowledgeRequestOption is a function that modifies a EventAcknowledgeRequest.
type EventAcknowledgeRequestOption func(*EventAcknowledgeRequest)

// WithSeverity sets the severity of the event.
// WithSeverity returns a EventAcknowledgeRequestOption that sets the severity.
func WithSeverity(severity Severity) EventAcknowledgeRequestOption {
	return func(e *EventAcknowledgeRequest) {
		e.Params.Severity = severity
	}
}

// WithActions sets the actions to perform on the event.
// WithActions returns a EventAcknowledgeRequestOption that sets the actions.
func WithActions(action ...EventAction) EventAcknowledgeRequestOption {
	actions := NewEventAction(action...)
	return func(e *EventAcknowledgeRequest) {
		e.Params.Action = actions
	}
}

// WithMessage sets the message to add to the event.
// WithMessage returns a EventAcknowledgeRequestOption that sets the message.
func WithMessage(message string) EventAcknowledgeRequestOption {
	return func(e *EventAcknowledgeRequest) {
		e.Params.Message = message
	}
}

// newEventAcknowledgeRequest creates a new event acknowledge request.
func newEventAcknowledgeRequest(eventids []string, opts ...EventAcknowledgeRequestOption) *EventAcknowledgeRequest {
	req := &EventAcknowledgeRequest{
		JSONRPC: JSONRPC,
		Method:  MethodEventAcknowledge,
		Params: EventsAcknowledgeParams{
			Eventids: eventids,
		},
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

// EventAcknowledgeResponse represents the response from an event acknowledgement.
type EventAcknowledgeResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Eventids []int `json:"eventids"`
	} `json:"result"`
	ID    int    `json:"id"`
	Error *Error `json:"error,omitempty"`
}

// AcknowledgeEvents acknowledges events with the specified options.
func (z *Client) AcknowledgeEvents(ctx context.Context, eventsID []string, opts ...EventAcknowledgeRequestOption) ([]int, error) {
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
	if res.Error != nil && res.Error.Code != 0 {
		return nil, res.Error
	}
	return res.Result.Eventids, nil
}

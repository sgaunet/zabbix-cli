package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const MethodProblemGet = "problem.get"

// Return only problems with given tags. Exact match by tag and case-insensitive search by value and operator.
// Format: [{"tag": "<tag>", "value": "<value>", "operator": "<operator>"}, ...].
// An empty array returns all problems.
// Possible operator types:
// 0 - (default) Like;
// 1 - Equal;
// 2 - Not like;
// 3 - Not equal
// 4 - Exists;
// 5 - Not exists.
type zbxTagsFilterProblem struct {
	Tag      string `json:"tag"      yaml:"tag"`
	Value    string `json:"value"    yaml:"value"`
	Operator string `json:"operator" yaml:"operator"`
}

// zbxParamsProblem represents the params for a problem.get request
type zbxParamsProblem struct {
	EventIDs     []string               `json:"eventids,omitempty"`
	GroupsIDs    []string               `json:"groupids,omitempty"`
	HostsIDs     []string               `json:"hostids,omitempty"`
	ObjectIDs    []string               `json:"objectids,omitempty"`
	Source       int                    `json:"source,omitempty"` // Return only problems with the given type. Refer to the problem event object page for a list of supported event types. Default: 0 - problem created by a trigger.
	Object       int                    `json:"object,omitempty"` // Return only problems with the given object type. Refer to the problem event object page for a list of supported object types. Default: 0 - trigger.
	Acknowledged bool                   `json:"acknowledged,omitempty"`
	Suppressed   bool                   `json:"suppressed,omitempty"`
	Severities   []string               `json:"severities,omitempty"` // Return only problems with given event severities. Applies only if object is trigger.
	EvalType     int                    `json:"evaltype,omitempty"`   // Rules for tag searching. Possible values: 0 - (default) And/Or;  2 - Or.
	Tags         []zbxTagsFilterProblem `json:"tags,omitempty"`
	Recent       bool                   `json:"recent,omitempty"`       // true - return PROBLEM and recently RESOLVED problems (depends on Display OK triggers for N seconds) Default: false - UNRESOLVED problems only
	EventidFrom  string                 `json:"eventid_from,omitempty"` // Return only problems with IDs greater or equal to the given ID.
	EventidTill  string                 `json:"eventid_till,omitempty"` // Return only problems with IDs less or equal to the given ID.
	TimeFrom     int64                  `json:"time_from,omitempty"`    // Return only problems that have been created after or at the given time.
	TimeTill     int64                  `json:"time_till,omitempty"`    // Return only problems that have been created before or at the given time.
	// selectAcknowledges 	query 	Return an acknowledges property with the problem updates. Problem updates are sorted in reverse chronological order.
	// The problem update object has the following properties:
	// acknowledgeid - (string) update's ID;
	// userid - (string) ID of the user that updated the event;
	// eventid - (string) ID of the updated event;
	// clock - (timestamp) time when the event was updated;
	// message - (string) text of the message;
	// action - (integer)type of update action (see event.acknowledge);
	// old_severity - (integer) event severity before this update action;
	// new_severity - (integer) event severity after this update action;

	// Supports count.
	// selectTags query // Return a tags property with the problem tags. Output format: [{"tag": "<tag>", "value": "<value>"}, ...].
	// selectSuppressionData 	query 	Return a suppression_data property with the list of maintenances:
	// maintenanceid - (string) ID of the maintenance;
	// suppress_until - (integer) time until the problem is suppressed.

	// sortfield 	string/array 	Sort the result by the given properties. Possible values are: eventid.
}

type zbxGetProblem struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  zbxParamsProblem `json:"params"`
	Auth    string           `json:"auth"`
	ID      int              `json:"id"`
}

type GetProblemOption func(*zbxGetProblem)

func GetProblemOptionEventIDs(eventIDs []string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.EventIDs = eventIDs
	}
}

func GetProblemOptionGroupsIDs(groupsIDs []string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.GroupsIDs = groupsIDs
	}
}

func GetProblemOptionHostsIDs(hostsIDs []string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.HostsIDs = hostsIDs
	}
}

func GetProblemOptionObjectIDs(objectIDs []string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.ObjectIDs = objectIDs
	}
}

func GetProblemOptionSource(source int) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Source = source
	}
}

func GetProblemOptionObject(object int) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Object = object
	}
}

func GetProblemOptionAcknowledged(acknowledged bool) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Acknowledged = acknowledged
	}
}

func GetProblemOptionSuppressed(suppressed bool) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Suppressed = suppressed
	}
}

func GetProblemOptionSeverities(severities []string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Severities = severities
	}
}

func GetProblemOptionEvalType(evalType int) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.EvalType = evalType
	}
}

func GetProblemOptionTags(tags []zbxTagsFilterProblem) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Tags = tags
	}
}

func GetProblemOptionRecent(recent bool) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.Recent = recent
	}
}

func GetProblemOptionEventidFrom(eventidFrom string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.EventidFrom = eventidFrom
	}
}

func GetProblemOptionEventidTill(eventidTill string) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.EventidTill = eventidTill
	}
}

func GetProblemOptionTimeFrom(timeFrom int64) GetProblemOption {
	return func(g *zbxGetProblem) {
		g.Params.TimeFrom = timeFrom
	}
}

// Problem represents a Zabbix problem
type Problem struct {
	Acknowledged  string `json:"acknowledged"`
	Clock         string `json:"clock"`
	CorrelationID string `json:"correlationID"`
	EventID       string `json:"eventID"`
	Name          string `json:"name"`
	Ns            string `json:"ns"`
	Object        string `json:"object"`
	ObjectID      string `json:"objectID"`
	Opdata        string `json:"opdata"`
	Rclock        string `json:"r_clock"`
	ReventID      string `json:"r_eventID"`
	Rns           string `json:"r_ns"`
	Severity      string `json:"severity"`
	Source        string `json:"source"`
	Suppressed    string `json:"suppressed"`
	UserID        string `json:"userid"`
}

// GetClock returns the clock as a time.Time
// If the clock cannot be converted, it returns time.Time{}
func (p *Problem) GetClock() time.Time {
	// convert clock to int64 and then to time.Time
	num, err := strconv.ParseInt(p.Clock, 10, 64)
	if err != nil {
		return time.Time{}
	}
	ts := time.Unix(num, 0)
	return ts
}

// GetRClock returns the rclock as a time.Time
// If the rclock cannot be converted or has no value, it returns time.Time{}
func (p *Problem) GetRClock() time.Time {
	if p.Rclock == "" {
		return time.Time{}
	}
	// convert clock to int64 and then to time.Time
	num, err := strconv.ParseInt(p.Rclock, 10, 64)
	if err != nil {
		return time.Time{}
	}
	ts := time.Unix(num, 0)
	return ts
}

func (p *Problem) GetDuration() time.Duration {
	emptyTime := time.Unix(0, 0)
	clock := p.GetClock()
	rclock := p.GetRClock()

	if clock == emptyTime {
		return time.Duration(0)
	}
	if rclock == emptyTime {
		return time.Since(clock)
	}
	return rclock.Sub(clock)
}

func (p *Problem) GetDurationStr() string {
	const NumberMinutesInHour = 60
	const NumberSecondsInMinute = 60
	durationProblem := p.GetDuration()
	return fmt.Sprintf("%d:%02d:%02d", int(durationProblem.Hours()), int(durationProblem.Minutes())%NumberMinutesInHour, int(durationProblem.Seconds())%NumberSecondsInMinute)
}

func (p *Problem) GetAcknowledge() bool {
	return p.Acknowledged == "1"
}

func (p *Problem) GetSuppressed() bool {
	return p.Suppressed == "1"
}

func (p *Problem) GetAcknowledgeStr() string {
	if p.GetAcknowledge() {
		return "Yes"
	}
	return "No"
}

func (p *Problem) GetSuppressedStr() string {
	if p.GetSuppressed() {
		return "Yes"
	}
	return "No"
}

func (p *Problem) GetSeverity() string {
	s, err := strconv.Atoi(p.Severity)
	if err != nil {
		return "Unknown"
	}
	return NewSeverity(s).String()
}

// zbxResultProblem represents the result of a problem.get request
type zbxResultProblem struct {
	JSONRPC  string    `json:"jsonrpc"`
	Result   []Problem `json:"result"`
	ErrorMsg ErrorMsg  `json:"error,omitempty"`
	ID       int       `json:"id"`
}

// GetProblems returns a list of problems
func (z *ZabbixAPI) GetProblems(ctx context.Context, opts ...GetProblemOption) ([]Problem, error) {
	payload := &zbxGetProblem{
		JSONRPC: JSONRPC,
		Method:  MethodProblemGet,
		Auth:    z.auth,
		ID:      generateUniqueID(),
	}
	payload.Auth = z.Auth()
	for _, opt := range opts {
		opt(payload)
	}

	statusCode, body, err := z.postRequest(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("cannot do request: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not OK: %d - %s (%w)", statusCode, string(body), ErrWrongHTTPCode)
	}

	var res zbxResultProblem
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w - %s", err, string(body))
	}
	if res.ErrorMsg != (ErrorMsg{}) {
		return nil, fmt.Errorf("error message: %w", &res.ErrorMsg)
	}
	return res.Result, nil
}

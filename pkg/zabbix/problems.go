package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// MethodProblemGet is the Zabbix API method for getting problems.
const MethodProblemGet = "problem.get"

// FilterProblemTags returns only problems with given tags. Exact match by tag and case-insensitive search by value and operator.
// Format: [{"tag": "<tag>", "value": "<value>", "operator": "<operator>"}, ...].
// An empty array returns all problems.
// Possible operator types:
// 0 - (default) Like;
// 1 - Equal;
// 2 - Not like;
// 3 - Not equal
// 4 - Exists;
// 5 - Not exists.
type FilterProblemTags struct {
	Tag      string `json:"tag"      yaml:"tag"`
	Value    string `json:"value"    yaml:"value"`
	Operator string `json:"operator" yaml:"operator"`
}

// ProblemParams represents the params for a problem.get request.
type ProblemParams struct {
	EventIDs     []string            `json:"eventids,omitempty"`
	GroupsIDs    []string            `json:"groupids,omitempty"`
	HostsIDs     []string            `json:"hostids,omitempty"`
	ObjectIDs    []string            `json:"objectids,omitempty"`
	Source       int                 `json:"source,omitempty"` // Return only problems with the given type. Refer to the problem event object page for a list of supported event types. Default: 0 - problem created by a trigger.
	Object       int                 `json:"object,omitempty"` // Return only problems with the given object type. Refer to the problem event object page for a list of supported object types. Default: 0 - trigger.
	Acknowledged bool                `json:"acknowledged"`
	Suppressed   bool                `json:"suppressed"`
	Severities   []string            `json:"severities,omitempty"` // Return only problems with given event severities. Applies only if object is trigger.
	EvalType     int                 `json:"evaltype,omitempty"`   // Rules for tag searching. Possible values: 0 - (default) And/Or;  2 - Or.
	Tags         []FilterProblemTags `json:"tags,omitempty"`
	Recent       bool                `json:"recent,omitempty"`       // true - return PROBLEM and recently RESOLVED problems (depends on Display OK triggers for N seconds) Default: false - UNRESOLVED problems only
	EventidFrom  string              `json:"eventid_from,omitempty"` // Return only problems with IDs greater or equal to the given ID.
	EventidTill  string              `json:"eventid_till,omitempty"` // Return only problems with IDs less or equal to the given ID.
	TimeFrom     int64               `json:"time_from,omitempty"`    // Return only problems that have been created after or at the given time.
	TimeTill     int64               `json:"time_till,omitempty"`    // Return only problems that have been created before or at the given time.
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
	CommonGetParams
}

// GetProblemRequest represents a request to retrieve problem information from Zabbix API.
type GetProblemRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  ProblemParams `json:"params"`
	Auth    string        `json:"auth"`
	ID      int           `json:"id"`
}

// GetProblemOption is a function that modifies a GetProblemRequest.
type GetProblemOption func(*GetProblemRequest)

// GetProblemOptionEventIDs sets event IDs as a filter option.
func GetProblemOptionEventIDs(eventIDs []string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.EventIDs = eventIDs
	}
}

// GetProblemOptionGroupsIDs sets group IDs as a filter option.
func GetProblemOptionGroupsIDs(groupsIDs []string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.GroupsIDs = groupsIDs
	}
}

// GetProblemOptionHostsIDs sets host IDs as a filter option.
func GetProblemOptionHostsIDs(hostsIDs []string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.HostsIDs = hostsIDs
	}
}

// GetProblemOptionObjectIDs sets object IDs as a filter option.
func GetProblemOptionObjectIDs(objectIDs []string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.ObjectIDs = objectIDs
	}
}

// GetProblemOptionSource sets the source as a filter option.
func GetProblemOptionSource(source int) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Source = source
	}
}

// GetProblemOptionObject sets the object as a filter option.
func GetProblemOptionObject(object int) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Object = object
	}
}

// GetProblemOptionAcknowledged sets the acknowledged state as a filter option.
func GetProblemOptionAcknowledged(acknowledged bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Acknowledged = acknowledged
	}
}

// GetProblemOptionSuppressed sets the suppressed state as a filter option.
func GetProblemOptionSuppressed(suppressed bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Suppressed = suppressed
	}
}

// GetProblemOptionSeverities sets the severities as a filter option.
func GetProblemOptionSeverities(severities []string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Severities = severities
	}
}

// GetProblemOptionEvalType sets the evaluation type as a filter option.
func GetProblemOptionEvalType(evalType int) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.EvalType = evalType
	}
}

// GetProblemOptionTags sets the tags as a filter option.
func GetProblemOptionTags(tags []FilterProblemTags) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Tags = tags
	}
}

// GetProblemOptionRecent sets the recent flag as a filter option.
func GetProblemOptionRecent(recent bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Recent = recent
	}
}

// GetProblemOptionEventidFrom sets the event ID from as a filter option.
func GetProblemOptionEventidFrom(eventidFrom string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.EventidFrom = eventidFrom
	}
}

// GetProblemOptionEventidTill sets the event ID till as a filter option.
func GetProblemOptionEventidTill(eventidTill string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.EventidTill = eventidTill
	}
}

// GetProblemOptionTimeFrom sets the time from as a filter option.
func GetProblemOptionTimeFrom(timeFrom int64) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.TimeFrom = timeFrom
	}
}

// GetProblemOptionTimeTill sets the time till as a filter option.
func GetProblemOptionTimeTill(timeTill int64) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.TimeTill = timeTill
	}
}

// GetProblemOptionCountOutput sets the count output flag as a filter option.
func GetProblemOptionCountOutput(countOutput bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.CountOutput = countOutput
	}
}

// GetProblemOptionEditable sets the editable flag as a filter option.
func GetProblemOptionEditable(editable bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Editable = editable
	}
}

// GetProblemOptionExcludeSearch sets the exclude search flag as a filter option.
func GetProblemOptionExcludeSearch(excludeSearch bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.ExcludeSearch = excludeSearch
	}
}

// GetProblemOptionLimit sets the limit as a filter option.
func GetProblemOptionLimit(limit int) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.Limit = limit
	}
}

// GetProblemOptionPreservekeys sets the preserve keys flag as a filter option.
func GetProblemOptionPreservekeys(preservekeys bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.PreserveKeys = preservekeys // Corrected field name to PreserveKeys
	}
}

// GetProblemOptionSearchByAny sets the search by any flag as a filter option.
func GetProblemOptionSearchByAny(searchByAny bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.SearchByAny = searchByAny
	}
}

// GetProblemOptionSearchWildcardsEnabled sets the search wildcards enabled flag as a filter option.
func GetProblemOptionSearchWildcardsEnabled(searchWildcardsEnabled bool) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.SearchWildcardsEnabled = searchWildcardsEnabled
	}
}

// Problem represents a Zabbix problem.
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

// GetClock returns the clock as a time.Time.
// If the clock cannot be converted, it returns time.Time{}.
func (p *Problem) GetClock() time.Time {
	// convert clock to int64 and then to time.Time
	num, err := strconv.ParseInt(p.Clock, 10, 64)
	if err != nil {
		return time.Time{}
	}
	ts := time.Unix(num, 0)
	return ts
}

// GetRClock returns the rclock as a time.Time.
// If the rclock cannot be converted or has no value, it returns time.Time{}.
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

// GetDuration returns the duration of the problem.
func (p *Problem) GetDuration() time.Duration {
	emptyTime := time.Unix(0, 0)
	clock := p.GetClock()
	rclock := p.GetRClock()

	if clock.Equal(emptyTime) {
		return time.Duration(0)
	}
	if rclock.Equal(emptyTime) {
		return time.Since(clock)
	}
	return rclock.Sub(clock)
}

// GetDurationStr returns the duration of the problem as a string.
func (p *Problem) GetDurationStr() string {
	const NumberMinutesInHour = 60
	const NumberSecondsInMinute = 60
	durationProblem := p.GetDuration()
	return fmt.Sprintf("%d:%02d:%02d", int(durationProblem.Hours()), int(durationProblem.Minutes())%NumberMinutesInHour, int(durationProblem.Seconds())%NumberSecondsInMinute)
}

// GetAcknowledge returns whether the problem is acknowledged.
func (p *Problem) GetAcknowledge() bool {
	return p.Acknowledged == "1"
}

// GetSuppressed returns whether the problem is suppressed.
func (p *Problem) GetSuppressed() bool {
	return p.Suppressed == "1"
}

// GetAcknowledgeStr returns the acknowledge state as a string.
func (p *Problem) GetAcknowledgeStr() string {
	if p.GetAcknowledge() {
		return "Yes"
	}
	return "No"
}

// GetSuppressedStr returns the suppressed state as a string.
func (p *Problem) GetSuppressedStr() string {
	if p.GetSuppressed() {
		return "Yes"
	}
	return "No"
}

// GetSeverity returns the severity of the problem.
func (p *Problem) GetSeverity() string {
	s, err := strconv.Atoi(p.Severity)
	if err != nil {
		return "Unknown"
	}
	return NewSeverity(s).String()
}

// ResultProblem represents the result of a problem.get request.
type ResultProblem struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  []Problem   `json:"result"`
	Error   *Error      `json:"error,omitempty"` // Changed to pointer type *Error
	ID      int         `json:"id"`
}

// GetProblems returns a list of problems.
func (z *Client) GetProblems(ctx context.Context, opts ...GetProblemOption) ([]Problem, error) {
	payload := &GetProblemRequest{
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

	var res ResultProblem
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w - %s", err, string(body))
	}
	if res.Error != nil && res.Error.Code != 0 {
		return nil, res.Error // Return the Error struct directly
	}
	return res.Result, nil
}

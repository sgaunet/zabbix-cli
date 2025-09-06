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
	CommonGetParams

	EventIDs     []string            `json:"eventids,omitempty"`
	GroupsIDs    []string            `json:"groupids,omitempty"`
	HostsIDs     []string            `json:"hostids,omitempty"`
	ObjectIDs    []string            `json:"objectids,omitempty"`
	Source       int                 `json:"source,omitempty"` // Default: 0 - problem created by a trigger.
	Object       int                 `json:"object,omitempty"` // Default: 0 - trigger.
	Acknowledged bool                `json:"acknowledged,omitempty"`
	Suppressed   bool                `json:"suppressed,omitempty"`
	Severities   []string            `json:"severities,omitempty"` // Applies only if object is trigger.
	EvalType     int                 `json:"evaltype,omitempty"`   // Rules for tag searching. 0 - (default) And/Or; 2 - Or.
	Tags         []FilterProblemTags `json:"tags,omitempty"`
	Recent       bool                `json:"recent,omitempty"`       // true - return PROBLEM and recently RESOLVED problems. Default: false - UNRESOLVED problems only.
	EventidFrom  string              `json:"eventid_from,omitempty"`
	EventidTill  string              `json:"eventid_till,omitempty"`
	TimeFrom     int64               `json:"time_from,omitempty"`
	TimeTill     int64               `json:"time_till,omitempty"`

	// Select parameters to include additional data in the response
	SelectAcknowledges    string `json:"selectAcknowledges,omitempty"`    // e.g., "extend", returns 'acknowledges' property in Problem.
	SelectTags            string `json:"selectTags,omitempty"`            // e.g., "extend", returns 'tags' property in Problem.
	SelectSuppressionData string `json:"selectSuppressionData,omitempty"` // e.g., "extend", returns 'suppression_data' property in Problem.
	SelectHosts           string `json:"selectHosts,omitempty"`           // To retrieve host information
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

// GetProblemOptionSelectHosts sets the selectHosts parameter for the problem.get request.
// Common values: "extend", or an array of specific properties like `["hostid", "name"]` (passed as a string).
func GetProblemOptionSelectHosts(selectQuery string) GetProblemOption {
	return func(g *GetProblemRequest) {
		g.Params.SelectHosts = selectQuery
	}
}

// ProblemResponseTag represents a tag associated with a problem, as returned by problem.get with selectTags.
// This is distinct from FilterProblemTags (used for filtering in params) and the ProblemTag in maintenance.go.
// API Reference: problem.get, selectTags parameter.
// Output format: [{"tag": "<tag>", "value": "<value>"}, ...]
type ProblemResponseTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// AcknowledgeEntry represents an acknowledgment record for a problem.
// API Reference: problem.get, selectAcknowledges parameter.
// The problem update object has the following properties:
// acknowledgeid - (string) update's ID;
// userid - (string) ID of the user that updated the event;
// eventid - (string) ID of the updated event;
// clock - (timestamp) time when the event was updated;
// message - (string) text of the message;
// action - (integer) type of update action (see event.acknowledge);
// old_severity - (integer) event severity before this update action;
// new_severity - (integer) event severity after this update action;
type AcknowledgeEntry struct {
	AcknowledgeID string `json:"acknowledgeid"`
	UserID        string `json:"userid"`
	EventID       string `json:"eventid"`
	Clock         string `json:"clock"`        // timestamp
	Message       string `json:"message"`
	Action        int    `json:"action"`       // integer
	OldSeverity   int    `json:"old_severity"` // integer
	NewSeverity   int    `json:"new_severity"` // integer
}

// SuppressionDataEntry represents data about problem suppression.
// API Reference: problem.get, selectSuppressionData parameter.
// The suppression_data object has the following properties:
// maintenanceid - (string) ID of the maintenance;
// suppress_until - (integer) time until the problem is suppressed.
type SuppressionDataEntry struct {
	MaintenanceID string `json:"maintenanceid"`
	SuppressUntil int64  `json:"suppress_until"` // timestamp
}

// HostInfo represents basic information about a host related to a problem.
// This is typically populated by a selectHosts query.
// Common fields are hostid and name.
type HostInfo struct {
	HostID string `json:"hostid"`
	Name   string `json:"name"`
	// Add other host properties here if needed, e.g., Host string `json:"host"` (technical name)
}

// Problem represents a Zabbix problem, potentially with additional data from select queries.
// API Reference: problem object, problem.get method.
// Most fields are Readonly.
type Problem struct {
	// Core problem fields
	EventID      string `json:"eventid"`                 // Readonly: ID of the problem event.
	Source       string `json:"source"`                  // Readonly: Type of object that created the problem event (e.g., "0" for trigger).
	Object       string `json:"object"`                  // Readonly: Type of object related to the problem event (e.g., "0" for trigger).
	ObjectID     string `json:"objectid"`                // Readonly: ID of the related object.
	Clock        string `json:"clock"`                   // Readonly: Time when the problem event was created (timestamp).
	Ns           string `json:"ns"`                      // Readonly: Nanoseconds when the problem event was created.
	Name         string `json:"name"`                    // Readonly: Problem name.
	Acknowledged string `json:"acknowledged"`            // Readonly: Whether the problem event is acknowledged ("0" or "1").
	Severity     string `json:"severity"`                // Readonly: Current severity of the problem (e.g., "0"-"5").

	// Recovery information
	Rclock   string `json:"r_clock,omitempty"`   // Readonly: Time when the problem was resolved (timestamp).
	ReventID string `json:"r_eventid,omitempty"` // Readonly: ID of the recovery event.
	Rns      string `json:"r_ns,omitempty"`      // Readonly: Nanoseconds when the problem was resolved.

	// Correlation information
	CorrelationID   string `json:"correlationid,omitempty"`   // Readonly: ID of the correlation rule that correlated the problem.
	CorrelationMode int    `json:"correlation_mode,omitempty"` // Readonly: How the problem was correlated.
	CorrelationTag  string `json:"correlation_tag,omitempty"`  // Readonly: Tag used for correlation.
	CauseEventID    string `json:"cause_eventid,omitempty"`    // Readonly: ID of the problem event that caused this problem.

	// Operational and suppression data
	Opdata        string `json:"opdata,omitempty"`           // Readonly: Operational data of the problem.
	Suppressed    string `json:"suppressed"`                 // Readonly: Whether the problem event is suppressed ("0" or "1").
	SuppressUntil string `json:"suppress_until,omitempty"`   // Readonly: Timestamp until when the problem event is suppressed (problem's own suppression time).

	// Fields populated by select queries
	Acknowledges    []AcknowledgeEntry     `json:"acknowledges,omitempty"`
	Tags            []ProblemResponseTag   `json:"tags,omitempty"` // Note: API returns "tags", distinct from ProblemParams.Tags used for filtering.
	SuppressionData []SuppressionDataEntry `json:"suppression_data,omitempty"`
	Hosts           []HostInfo             `json:"hosts,omitempty"` // Populated by selectHosts
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
	JSONRPC string        `json:"jsonrpc"`
	Result  []Problem     `json:"result"`
	Error   *Error        `json:"error,omitempty"` // Uses common.Error as per MEMORY[cf628beb-6a65-471f-8fc1-8c971bbf63fe] (defined in pkg/zabbix/common.go)
	ID      int           `json:"id"`
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

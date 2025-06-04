package zabbix

// HostGroupGetParams defines the parameters for the Zabbix hostgroup.get API call.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/get
type HostGroupGetParams struct {
	CommonGetParams // Embeds common parameters like Output, Limit, Filter, etc.

	GroupIDs                []string `json:"groupids,omitempty"`
	GraphIDs                []string `json:"graphids,omitempty"`
	HostIDs                 []string `json:"hostids,omitempty"`
	MaintenanceIDs          []string `json:"maintenanceids,omitempty"`
	MonitoredHosts          bool     `json:"monitored_hosts,omitempty"`          // Return only host groups that contain monitored hosts.
	RealHosts               bool     `json:"real_hosts,omitempty"`               // Return only host groups that contain hosts.
	TemplatedHosts          bool     `json:"templated_hosts,omitempty"`          // Return only host groups that contain templates.
	WithApplications        bool     `json:"with_applications,omitempty"`        // Return only host groups that contain applications.
	WithGraphs              bool     `json:"with_graphs,omitempty"`              // Return only host groups that contain graphs.
	WithHistoricalItems     bool     `json:"with_historical_items,omitempty"`    // Return only host groups that contain items with history data.
	WithHostsAndTemplates   bool     `json:"with_hosts_and_templates,omitempty"` // Return only host groups that contain hosts or templates.
	WithHTTPTests           bool     `json:"with_httptests,omitempty"`           // Return only host groups that contain web scenarios.
	WithItems               bool     `json:"with_items,omitempty"`               // Return only host groups that contain items.
	WithMonitoredHTTPTests  bool     `json:"with_monitored_httptests,omitempty"` // Return only host groups that contain enabled web scenarios that belong to monitored hosts.
	WithMonitoredItems      bool     `json:"with_monitored_items,omitempty"`     // Return only host groups that contain enabled items that belong to monitored hosts.
	WithMonitoredTriggers   bool     `json:"with_monitored_triggers,omitempty"`  // Return only host groups that contain enabled triggers that belong to monitored hosts and are not in a problem state.
	WithSimpleGraphItems    bool     `json:"with_simple_graph_items,omitempty"`  // Return only host groups that contain items used in simple graphs.
	WithTriggers            bool     `json:"with_triggers,omitempty"`            // Return only host groups that contain triggers.
	Editable                bool     `json:"editable,omitempty"`                 // Return only host groups that are available for writing.

	SelectApplications      interface{} `json:"selectApplications,omitempty"`      // "extend" or array of fields
	SelectGroupDiscovery    interface{} `json:"selectGroupDiscovery,omitempty"`    // "extend" or array of fields
	SelectHosts             interface{} `json:"selectHosts,omitempty"`             // "extend" or array of fields
	SelectHTTPTests         interface{} `json:"selectHttptests,omitempty"`         // "extend" or array of fields
	SelectItems             interface{} `json:"selectItems,omitempty"`             // "extend" or array of fields
	SelectMonitoredHosts    interface{} `json:"selectMonitoredHosts,omitempty"`    // "extend" or array of fields
	SelectRealHosts         interface{} `json:"selectRealHosts,omitempty"`         // "extend" or array of fields
	SelectSimpleGraphItems  interface{} `json:"selectSimpleGraphItems,omitempty"`  // "extend" or array of fields
	SelectTriggers          interface{} `json:"selectTriggers,omitempty"`          // "extend" or array of fields
	LimitSelects            int         `json:"limitSelects,omitempty"`            // Limits the number of records returned by subselects.
}

// HostGroupGetRequest defines the JSON-RPC request structure for hostgroup.get.
type HostGroupGetRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  HostGroupGetParams `json:"params"`
	Auth    string             `json:"auth,omitempty"`
	ID      int                `json:"id"`
}

// HostGroupGetResponse defines the JSON-RPC response structure for hostgroup.get.
type HostGroupGetResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  []HostGroup  `json:"result"` // Array of HostGroup objects
	ID      int          `json:"id"`
	Error   *Error       `json:"error,omitempty"`
}

// HostGroupGetOption defines a function signature for options to configure a HostGroupGetRequest.
type HostGroupGetOption func(*HostGroupGetRequest)

// NewHostGroupGetRequest creates a new HostGroupGetRequest with default values and applies any provided options.
func NewHostGroupGetRequest(options ...HostGroupGetOption) *HostGroupGetRequest {
	hgr := &HostGroupGetRequest{
		JSONRPC: "2.0",
		Method:  "hostgroup.get",
		Params:  HostGroupGetParams{}, // Initialize with empty params
	}
	for _, opt := range options {
		opt(hgr)
	}
	return hgr
}

// --- Option functions for HostGroupGetParams ---

// WithHostGroupGetGroupIDs sets the group IDs for the request.
func WithHostGroupGetGroupIDs(ids []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.GroupIDs = ids }
}

// WithHostGroupGetHostIDs sets the host IDs for the request.
func WithHostGroupGetHostIDs(ids []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.HostIDs = ids }
}

// WithHostGroupGetGraphIDs sets the graph IDs for the request.
func WithHostGroupGetGraphIDs(ids []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.GraphIDs = ids }
}

// WithHostGroupGetMaintenanceIDs sets the maintenance IDs for the request.
func WithHostGroupGetMaintenanceIDs(ids []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.MaintenanceIDs = ids }
}

// WithHostGroupGetMonitoredHosts sets the monitored_hosts flag.
func WithHostGroupGetMonitoredHosts(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.MonitoredHosts = flag }
}

// WithHostGroupGetRealHosts sets the real_hosts flag.
func WithHostGroupGetRealHosts(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.RealHosts = flag }
}

// WithHostGroupGetTemplatedHosts sets the templated_hosts flag.
func WithHostGroupGetTemplatedHosts(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.TemplatedHosts = flag }
}

// WithHostGroupGetWithApplications sets the with_applications flag.
func WithHostGroupGetWithApplications(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithApplications = flag }
}

// WithHostGroupGetWithGraphs sets the with_graphs flag.
func WithHostGroupGetWithGraphs(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithGraphs = flag }
}

// WithHostGroupGetWithHistoricalItems sets the with_historical_items flag.
func WithHostGroupGetWithHistoricalItems(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithHistoricalItems = flag }
}

// WithHostGroupGetWithHostsAndTemplates sets the with_hosts_and_templates flag.
func WithHostGroupGetWithHostsAndTemplates(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithHostsAndTemplates = flag }
}

// WithHostGroupGetWithHTTPTests sets the with_httptests flag.
func WithHostGroupGetWithHTTPTests(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithHTTPTests = flag }
}

// WithHostGroupGetWithItems sets the with_items flag.
func WithHostGroupGetWithItems(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithItems = flag }
}

// WithHostGroupGetWithMonitoredHTTPTests sets the with_monitored_httptests flag.
func WithHostGroupGetWithMonitoredHTTPTests(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithMonitoredHTTPTests = flag }
}

// WithHostGroupGetWithMonitoredItems sets the with_monitored_items flag.
func WithHostGroupGetWithMonitoredItems(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithMonitoredItems = flag }
}

// WithHostGroupGetWithMonitoredTriggers sets the with_monitored_triggers flag.
func WithHostGroupGetWithMonitoredTriggers(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithMonitoredTriggers = flag }
}

// WithHostGroupGetWithSimpleGraphItems sets the with_simple_graph_items flag.
func WithHostGroupGetWithSimpleGraphItems(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithSimpleGraphItems = flag }
}

// WithHostGroupGetWithTriggers sets the with_triggers flag.
func WithHostGroupGetWithTriggers(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.WithTriggers = flag }
}

// WithHostGroupGetEditable sets the editable flag.
func WithHostGroupGetEditable(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.Editable = flag }
}

// WithHostGroupGetSelectApplications sets the selectApplications parameter.
func WithHostGroupGetSelectApplications(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectApplications = query }
}

// WithHostGroupGetSelectGroupDiscovery sets the selectGroupDiscovery parameter.
func WithHostGroupGetSelectGroupDiscovery(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectGroupDiscovery = query }
}

// WithHostGroupGetSelectHosts sets the selectHosts parameter.
func WithHostGroupGetSelectHosts(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectHosts = query }
}

// WithHostGroupGetSelectHTTPTests sets the selectHttptests parameter.
func WithHostGroupGetSelectHTTPTests(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectHTTPTests = query }
}

// WithHostGroupGetSelectItems sets the selectItems parameter.
func WithHostGroupGetSelectItems(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectItems = query }
}

// WithHostGroupGetSelectMonitoredHosts sets the selectMonitoredHosts parameter.
func WithHostGroupGetSelectMonitoredHosts(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectMonitoredHosts = query }
}

// WithHostGroupGetSelectRealHosts sets the selectRealHosts parameter.
func WithHostGroupGetSelectRealHosts(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectRealHosts = query }
}

// WithHostGroupGetSelectSimpleGraphItems sets the selectSimpleGraphItems parameter.
func WithHostGroupGetSelectSimpleGraphItems(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectSimpleGraphItems = query }
}

// WithHostGroupGetSelectTriggers sets the selectTriggers parameter.
func WithHostGroupGetSelectTriggers(query interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SelectTriggers = query }
}

// WithHostGroupGetLimitSelects sets the limitSelects parameter.
func WithHostGroupGetLimitSelects(limit int) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.LimitSelects = limit }
}

// --- Option functions for CommonGetParams (embedded) ---

// WithHostGroupGetOutput sets the output parameter.
func WithHostGroupGetOutput(output interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.Output = output }
}

// WithHostGroupGetLimit sets the limit parameter.
func WithHostGroupGetLimit(limit int) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.Limit = limit }
}

// WithHostGroupGetFilter sets the filter parameter.
func WithHostGroupGetFilter(filter map[string]interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.Filter = filter }
}

// WithHostGroupGetSortField sets the sortfield parameter.
func WithHostGroupGetSortField(sortField []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SortField = sortField }
}

// WithHostGroupGetSortOrder sets the sortorder parameter.
func WithHostGroupGetSortOrder(sortOrder []string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SortOrder = sortOrder }
}

// WithHostGroupGetSearch sets the search parameter.
func WithHostGroupGetSearch(search map[string]interface{}) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.Search = search }
}

// WithHostGroupGetSearchByAny sets the searchByAny flag.
func WithHostGroupGetSearchByAny(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SearchByAny = flag }
}

// WithHostGroupGetSearchWildcardsEnabled sets the searchWildcardsEnabled flag.
func WithHostGroupGetSearchWildcardsEnabled(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.SearchWildcardsEnabled = flag }
}

// WithHostGroupGetStartSearch sets the startSearch flag.
func WithHostGroupGetStartSearch(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.StartSearch = flag }
}

// WithHostGroupGetExcludeSearch sets the excludeSearch flag.
func WithHostGroupGetExcludeSearch(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.ExcludeSearch = flag }
}

// WithHostGroupGetPreserveKeys sets the preservekeys flag.
func WithHostGroupGetPreserveKeys(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.PreserveKeys = flag }
}

// WithHostGroupGetCountOutput sets the countOutput flag.
func WithHostGroupGetCountOutput(flag bool) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Params.CountOutput = flag }
}

// --- Option functions for request Auth and ID ---

// WithHostGroupGetAuth sets the authentication token for the API request.
func WithHostGroupGetAuth(token string) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.Auth = token }
}

// WithHostGroupGetID sets the ID for the API request.
func WithHostGroupGetID(id int) HostGroupGetOption {
	return func(hgr *HostGroupGetRequest) { hgr.ID = id }
}

// NewGetAllHostGroupsRequest creates a HostGroupGetRequest configured to fetch all host groups
// with all their properties. It allows overriding default auth and ID via options.
func NewGetAllHostGroupsRequest(options ...HostGroupGetOption) *HostGroupGetRequest {
	// Define the base options for getting all host groups
	baseOptions := []HostGroupGetOption{
		WithHostGroupGetOutput("extend"), // Ensure all fields are fetched for host groups
	}

	// Combine base options with any user-provided options.
	// User-provided options can override base options if they modify the same fields.
	// User options like Auth and ID will be applied.
	allOptions := append(baseOptions, options...)

	// Use the existing NewHostGroupGetRequest constructor with the combined options
	return NewHostGroupGetRequest(allOptions...)
}


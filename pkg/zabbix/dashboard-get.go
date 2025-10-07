package zabbix

// DashboardGetParams defines the parameters for the Zabbix dashboard.get API call.
// See: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/dashboard/get
type DashboardGetParams struct {
	CommonGetParams // Embeds common parameters like Output, Limit, Filter, etc.

	DashboardIDs []string `json:"dashboardids,omitempty"` // Return only dashboards with the given IDs.

	SelectPages      interface{} `json:"selectPages,omitempty"`      // Return dashboard pages. "extend" or array of fields
	SelectUsers      interface{} `json:"selectUsers,omitempty"`      // Return users the dashboard is shared with. "extend" or array of fields
	SelectUserGroups interface{} `json:"selectUserGroups,omitempty"` // Return user groups the dashboard is shared with. "extend" or array of fields
}

// DashboardGetRequest defines the JSON-RPC request structure for dashboard.get.
type DashboardGetRequest struct {
	JSONRPC string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  DashboardGetParams  `json:"params"`
	Auth    string              `json:"auth,omitempty"`
	ID      int                 `json:"id"`
}

// DashboardGetResponse defines the JSON-RPC response structure for dashboard.get.
type DashboardGetResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  []Dashboard `json:"result"` // Array of Dashboard objects
	ID      int         `json:"id"`
	Error   *Error      `json:"error,omitempty"`
}

// DashboardGetOption defines a function signature for options to configure a DashboardGetRequest.
type DashboardGetOption func(*DashboardGetRequest)

// NewDashboardGetRequest creates a new DashboardGetRequest with default values and applies any provided options.
func NewDashboardGetRequest(options ...DashboardGetOption) *DashboardGetRequest {
	dgr := &DashboardGetRequest{
		JSONRPC: JSONRPC,
		Method:  "dashboard.get",
		Params:  DashboardGetParams{}, // Initialize with empty params
	}
	for _, opt := range options {
		opt(dgr)
	}
	return dgr
}

// --- Option functions for DashboardGetParams ---

// WithDashboardGetDashboardIDs sets the dashboard IDs for the request.
func WithDashboardGetDashboardIDs(ids []string) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.DashboardIDs = ids }
}

// WithDashboardGetSelectPages sets the selectPages parameter.
func WithDashboardGetSelectPages(query interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SelectPages = query }
}

// WithDashboardGetSelectUsers sets the selectUsers parameter.
func WithDashboardGetSelectUsers(query interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SelectUsers = query }
}

// WithDashboardGetSelectUserGroups sets the selectUserGroups parameter.
func WithDashboardGetSelectUserGroups(query interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SelectUserGroups = query }
}

// --- Option functions for CommonGetParams (embedded) ---

// WithDashboardGetOutput sets the output parameter.
func WithDashboardGetOutput(output interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.Output = output }
}

// WithDashboardGetLimit sets the limit parameter.
func WithDashboardGetLimit(limit int) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.Limit = limit }
}

// WithDashboardGetFilter sets the filter parameter.
func WithDashboardGetFilter(filter map[string]interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.Filter = filter }
}

// WithDashboardGetSortField sets the sortfield parameter.
func WithDashboardGetSortField(sortField []string) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SortField = sortField }
}

// WithDashboardGetSortOrder sets the sortorder parameter.
func WithDashboardGetSortOrder(sortOrder []string) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SortOrder = sortOrder }
}

// WithDashboardGetSearch sets the search parameter.
func WithDashboardGetSearch(search map[string]interface{}) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.Search = search }
}

// WithDashboardGetSearchByAny sets the searchByAny flag.
func WithDashboardGetSearchByAny(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SearchByAny = flag }
}

// WithDashboardGetSearchWildcardsEnabled sets the searchWildcardsEnabled flag.
func WithDashboardGetSearchWildcardsEnabled(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.SearchWildcardsEnabled = flag }
}

// WithDashboardGetStartSearch sets the startSearch flag.
func WithDashboardGetStartSearch(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.StartSearch = flag }
}

// WithDashboardGetExcludeSearch sets the excludeSearch flag.
func WithDashboardGetExcludeSearch(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.ExcludeSearch = flag }
}

// WithDashboardGetPreserveKeys sets the preservekeys flag.
func WithDashboardGetPreserveKeys(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.PreserveKeys = flag }
}

// WithDashboardGetCountOutput sets the countOutput flag.
func WithDashboardGetCountOutput(flag bool) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Params.CountOutput = flag }
}

// --- Option functions for request Auth and ID ---

// WithDashboardGetAuth sets the authentication token for the API request.
func WithDashboardGetAuth(token string) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.Auth = token }
}

// WithDashboardGetID sets the ID for the API request.
func WithDashboardGetID(id int) DashboardGetOption {
	return func(dgr *DashboardGetRequest) { dgr.ID = id }
}

// NewGetAllDashboardsRequest creates a DashboardGetRequest configured to fetch all dashboards
// with all their properties and pages. It allows overriding default auth and ID via options.
func NewGetAllDashboardsRequest(options ...DashboardGetOption) *DashboardGetRequest {
	// Define the base options for getting all dashboards
	baseOptions := []DashboardGetOption{
		WithDashboardGetOutput("extend"),       // Ensure all fields are fetched for dashboards
		WithDashboardGetSelectPages("extend"),  // Fetch all pages with widgets
	}

	// Create a new slice with enough capacity to hold both base and user options
	allOptions := make([]DashboardGetOption, 0, len(baseOptions)+len(options))

	// Add base options first
	allOptions = append(allOptions, baseOptions...)
	// Then add user-provided options which can override base options
	allOptions = append(allOptions, options...)

	// Use the existing NewDashboardGetRequest constructor with the combined options
	return NewDashboardGetRequest(allOptions...)
}

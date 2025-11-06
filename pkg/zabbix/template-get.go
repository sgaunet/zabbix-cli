package zabbix

// TemplateGetParams defines the parameters for the Zabbix template.get API call.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/template/get
type TemplateGetParams struct {
	CommonGetParams // Embeds common parameters like Output, Limit, Filter, etc.

	TemplateIDs        []string `json:"templateids,omitempty"`        // Return only templates with the given template IDs.
	GroupIDs           []string `json:"groupids,omitempty"`           // Return only templates that belong to the given groups.
	ParentTemplateIDs  []string `json:"parentTemplateids,omitempty"`  // Return only templates that are children of the given templates.
	HostIDs            []string `json:"hostids,omitempty"`            // Return only templates that are linked to the given hosts.
	GraphIDs           []string `json:"graphids,omitempty"`           // Return only templates that contain the given graphs.
	ItemIDs            []string `json:"itemids,omitempty"`            // Return only templates that contain the given items.
	TriggerIDs         []string `json:"triggerids,omitempty"`         // Return only templates that contain the given triggers.

	WithItems          bool     `json:"with_items,omitempty"`         // Return only templates that have items.
	WithTriggers       bool     `json:"with_triggers,omitempty"`      // Return only templates that have triggers.
	WithGraphs         bool     `json:"with_graphs,omitempty"`        // Return only templates that have graphs.
	WithHTTPTests      bool     `json:"with_httptests,omitempty"`     // Return only templates that have web scenarios.

	SelectTags              any `json:"selectTags,omitempty"`              // "extend" or array of fields
	SelectHosts             any `json:"selectHosts,omitempty"`             // "extend" or array of fields
	SelectTemplateGroups    any `json:"selectTemplateGroups,omitempty"`    // "extend" or array of fields (formerly selectGroups in older versions)
	SelectTemplates         any `json:"selectTemplates,omitempty"`         // "extend" or array of fields (child templates)
	SelectParentTemplates   any `json:"selectParentTemplates,omitempty"`   // "extend" or array of fields
	SelectHTTPTests         any `json:"selectHttpTests,omitempty"`         // "extend" or array of fields
	SelectItems             any `json:"selectItems,omitempty"`             // "extend" or array of fields
	SelectDiscoveries       any `json:"selectDiscoveries,omitempty"`       // "extend" or array of fields
	SelectTriggers          any `json:"selectTriggers,omitempty"`          // "extend" or array of fields
	SelectGraphs            any `json:"selectGraphs,omitempty"`            // "extend" or array of fields
	SelectMacros            any `json:"selectMacros,omitempty"`            // "extend" or array of fields
	SelectDashboards        any `json:"selectDashboards,omitempty"`        // "extend" or array of fields
	SelectValueMaps         any `json:"selectValueMaps,omitempty"`         // "extend" or array of fields
}

// TemplateGetRequest defines the JSON-RPC request structure for template.get.
type TemplateGetRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  TemplateGetParams  `json:"params"`
	Auth    string             `json:"auth,omitempty"`
	ID      int                `json:"id"`
}

// Template represents a template object returned by template.get.
// Based on Zabbix API 7.2 specification.
type Template struct {
	TemplateID  string `json:"templateid"`            // ID of the template.
	Host        string `json:"host"`                  // Technical name of the template.
	Name        string `json:"name,omitempty"`        // Visible name of the template.
	Description string `json:"description,omitempty"` // Description of the template.
	UUID        string `json:"uuid,omitempty"`        // Universal unique identifier for template linking.
	Vendor      struct {
		Name    string `json:"name,omitempty"`    // Vendor name.
		Version string `json:"version,omitempty"` // Vendor version.
	} `json:"vendor,omitempty"` // Template vendor information.
}

// TemplateGetResponse defines the JSON-RPC response structure for template.get.
type TemplateGetResponse struct {
	JSONRPC string     `json:"jsonrpc"`
	Result  []Template `json:"result"` // Array of Template objects
	ID      int        `json:"id"`
	Error   *Error     `json:"error,omitempty"`
}

// TemplateGetOption defines a function signature for options to configure a TemplateGetRequest.
type TemplateGetOption func(*TemplateGetRequest)

// NewTemplateGetRequest creates a new TemplateGetRequest with default values and applies any provided options.
func NewTemplateGetRequest(options ...TemplateGetOption) *TemplateGetRequest {
	tgr := &TemplateGetRequest{
		JSONRPC: JSONRPC,
		Method:  "template.get",
		Params:  TemplateGetParams{}, // Initialize with empty params
		ID:      generateUniqueID(),
	}
	for _, opt := range options {
		opt(tgr)
	}
	return tgr
}

// --- Option functions for TemplateGetParams ---

// WithTemplateGetTemplateIDs sets the template IDs for the request.
func WithTemplateGetTemplateIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.TemplateIDs = ids }
}

// WithTemplateGetGroupIDs sets the group IDs for the request.
func WithTemplateGetGroupIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.GroupIDs = ids }
}

// WithTemplateGetParentTemplateIDs sets the parent template IDs for the request.
func WithTemplateGetParentTemplateIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.ParentTemplateIDs = ids }
}

// WithTemplateGetHostIDs sets the host IDs for the request.
func WithTemplateGetHostIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.HostIDs = ids }
}

// WithTemplateGetGraphIDs sets the graph IDs for the request.
func WithTemplateGetGraphIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.GraphIDs = ids }
}

// WithTemplateGetItemIDs sets the item IDs for the request.
func WithTemplateGetItemIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.ItemIDs = ids }
}

// WithTemplateGetTriggerIDs sets the trigger IDs for the request.
func WithTemplateGetTriggerIDs(ids []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.TriggerIDs = ids }
}

// WithTemplateGetWithItems sets the with_items flag.
func WithTemplateGetWithItems(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.WithItems = flag }
}

// WithTemplateGetWithTriggers sets the with_triggers flag.
func WithTemplateGetWithTriggers(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.WithTriggers = flag }
}

// WithTemplateGetWithGraphs sets the with_graphs flag.
func WithTemplateGetWithGraphs(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.WithGraphs = flag }
}

// WithTemplateGetWithHTTPTests sets the with_httptests flag.
func WithTemplateGetWithHTTPTests(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.WithHTTPTests = flag }
}

// WithTemplateGetSelectTags sets the selectTags parameter.
func WithTemplateGetSelectTags(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectTags = query }
}

// WithTemplateGetSelectHosts sets the selectHosts parameter.
func WithTemplateGetSelectHosts(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectHosts = query }
}

// WithTemplateGetSelectTemplateGroups sets the selectTemplateGroups parameter.
func WithTemplateGetSelectTemplateGroups(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectTemplateGroups = query }
}

// WithTemplateGetSelectTemplates sets the selectTemplates parameter (child templates).
func WithTemplateGetSelectTemplates(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectTemplates = query }
}

// WithTemplateGetSelectParentTemplates sets the selectParentTemplates parameter.
func WithTemplateGetSelectParentTemplates(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectParentTemplates = query }
}

// WithTemplateGetSelectHTTPTests sets the selectHttpTests parameter.
func WithTemplateGetSelectHTTPTests(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectHTTPTests = query }
}

// WithTemplateGetSelectItems sets the selectItems parameter.
func WithTemplateGetSelectItems(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectItems = query }
}

// WithTemplateGetSelectDiscoveries sets the selectDiscoveries parameter.
func WithTemplateGetSelectDiscoveries(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectDiscoveries = query }
}

// WithTemplateGetSelectTriggers sets the selectTriggers parameter.
func WithTemplateGetSelectTriggers(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectTriggers = query }
}

// WithTemplateGetSelectGraphs sets the selectGraphs parameter.
func WithTemplateGetSelectGraphs(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectGraphs = query }
}

// WithTemplateGetSelectMacros sets the selectMacros parameter.
func WithTemplateGetSelectMacros(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectMacros = query }
}

// WithTemplateGetSelectDashboards sets the selectDashboards parameter.
func WithTemplateGetSelectDashboards(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectDashboards = query }
}

// WithTemplateGetSelectValueMaps sets the selectValueMaps parameter.
func WithTemplateGetSelectValueMaps(query any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SelectValueMaps = query }
}

// --- Option functions for CommonGetParams (embedded) ---

// WithTemplateGetOutput sets the output parameter.
func WithTemplateGetOutput(output any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.Output = output }
}

// WithTemplateGetLimit sets the limit parameter.
func WithTemplateGetLimit(limit int) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.Limit = limit }
}

// WithTemplateGetFilter sets the filter parameter.
func WithTemplateGetFilter(filter map[string]any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.Filter = filter }
}

// WithTemplateGetSortField sets the sortfield parameter.
func WithTemplateGetSortField(sortField []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SortField = sortField }
}

// WithTemplateGetSortOrder sets the sortorder parameter.
func WithTemplateGetSortOrder(sortOrder []string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SortOrder = sortOrder }
}

// WithTemplateGetSearch sets the search parameter.
func WithTemplateGetSearch(search map[string]any) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.Search = search }
}

// WithTemplateGetSearchByAny sets the searchByAny flag.
func WithTemplateGetSearchByAny(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SearchByAny = flag }
}

// WithTemplateGetSearchWildcardsEnabled sets the searchWildcardsEnabled flag.
func WithTemplateGetSearchWildcardsEnabled(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.SearchWildcardsEnabled = flag }
}

// WithTemplateGetStartSearch sets the startSearch flag.
func WithTemplateGetStartSearch(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.StartSearch = flag }
}

// WithTemplateGetExcludeSearch sets the excludeSearch flag.
func WithTemplateGetExcludeSearch(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.ExcludeSearch = flag }
}

// WithTemplateGetPreserveKeys sets the preservekeys flag.
func WithTemplateGetPreserveKeys(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.PreserveKeys = flag }
}

// WithTemplateGetCountOutput sets the countOutput flag.
func WithTemplateGetCountOutput(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.CountOutput = flag }
}

// WithTemplateGetEditable sets the editable flag.
func WithTemplateGetEditable(flag bool) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Params.Editable = flag }
}

// --- Option functions for request Auth and ID ---

// WithTemplateGetAuth sets the authentication token for the API request.
func WithTemplateGetAuth(token string) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.Auth = token }
}

// WithTemplateGetID sets the ID for the API request.
func WithTemplateGetID(id int) TemplateGetOption {
	return func(tgr *TemplateGetRequest) { tgr.ID = id }
}

// NewGetAllTemplatesRequest creates a TemplateGetRequest configured to fetch all templates
// with all their properties. It allows overriding default auth and ID via options.
func NewGetAllTemplatesRequest(options ...TemplateGetOption) *TemplateGetRequest {
	// Define the base options for getting all templates
	baseOptions := []TemplateGetOption{
		WithTemplateGetOutput("extend"), // Ensure all fields are fetched for templates
	}

	// Create a new slice with enough capacity to hold both base and user options
	allOptions := make([]TemplateGetOption, 0, len(baseOptions) + len(options))

	// Add base options first
	allOptions = append(allOptions, baseOptions...)
	// Then add user-provided options which can override base options
	allOptions = append(allOptions, options...)

	// Use the existing NewTemplateGetRequest constructor with the combined options
	return NewTemplateGetRequest(allOptions...)
}

// GetTemplateID returns the IDs of the templates from the response.
func (t *TemplateGetResponse) GetTemplateID() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.TemplateID)
	}
	return templates
}

// GetTemplateName returns the names of the templates from the response.
func (t *TemplateGetResponse) GetTemplateName() []string {
	var templates []string //nolint: prealloc
	for _, template := range t.Result {
		templates = append(templates, template.Name)
	}
	return templates
}

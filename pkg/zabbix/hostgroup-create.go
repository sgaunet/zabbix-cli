package zabbix

// HostGroupCreateRequest represents the JSON-RPC request for hostgroup.create.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/create
type HostGroupCreateRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  []HostGroup `json:"params"` // Array of HostGroup objects to create. Only 'name' is required.
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

// HostGroupCreateResponseData contains the 'groupids' from the hostgroup.create response.
type HostGroupCreateResponseData struct {
	GroupIDs []string `json:"groupids"`
}

// HostGroupCreateResponse represents the JSON-RPC response for hostgroup.create.
type HostGroupCreateResponse struct {
	JSONRPC string                      `json:"jsonrpc"`
	Result  HostGroupCreateResponseData `json:"result"`
	ID      int                         `json:"id"`
	Error   *Error                      `json:"error,omitempty"` // Assuming an Error struct exists
}

// HostGroupCreateOption defines a function signature for options to configure a HostGroupCreateRequest.
type HostGroupCreateOption func(*HostGroupCreateRequest)

// NewHostGroupCreateRequest creates a new HostGroupCreateRequest with default values and applies any provided options.
// It takes one or more host group names to create.
func NewHostGroupCreateRequest(groupNames []string, options ...HostGroupCreateOption) *HostGroupCreateRequest {
	if len(groupNames) == 0 {
		// Or handle error: return nil, errors.New("at least one group name must be provided")
		return nil
	}
	params := make([]HostGroup, len(groupNames))
	for i, name := range groupNames {
		params[i] = HostGroup{Name: name}
	}

	hcr := &HostGroupCreateRequest{
		JSONRPC: JSONRPC,
		Method:  "hostgroup.create",
		Params:  params,
	}

	for _, opt := range options {
		opt(hcr)
	}
	return hcr
}

// WithHostGroupCreateAuth sets the authentication token for the API request.
func WithHostGroupCreateAuth(token string) HostGroupCreateOption {
	return func(hcr *HostGroupCreateRequest) {
		hcr.Auth = token
	}
}

// WithHostGroupCreateID sets the ID for the API request.
func WithHostGroupCreateID(id int) HostGroupCreateOption {
	return func(hcr *HostGroupCreateRequest) {
		hcr.ID = id
	}
}

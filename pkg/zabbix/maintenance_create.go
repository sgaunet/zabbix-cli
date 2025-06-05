package zabbix

// MaintenanceCreateOption defines a function signature for options to configure a MaintenanceCreateRequest.
// These options are used with NewMaintenanceCreateRequest to customize the request.
type MaintenanceCreateOption func(*MaintenanceCreateRequest)

// NewMaintenanceCreateRequest creates a new MaintenanceCreateRequest with default values and applies any provided options.
// Default JSONRPC version is "2.0" and method is "maintenance.create".
func NewMaintenanceCreateRequest(options ...MaintenanceCreateOption) *MaintenanceCreateRequest {
	mcr := &MaintenanceCreateRequest{
		JSONRPC: JSONRPC,
		Method:  "maintenance.create",
		Params:  Maintenance{}, // Initialize empty Maintenance params
	}
	for _, opt := range options {
		opt(mcr)
	}
	return mcr
}

// WithMaintenanceName sets the name for the maintenance period.
func WithMaintenanceName(name string) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.Name = name
	}
}

// WithMaintenanceActiveSince sets the time when the maintenance period becomes active (Unix timestamp).
func WithMaintenanceActiveSince(activeSince int64) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.ActiveSince = StringInt64(activeSince)
	}
}

// WithMaintenanceActiveTill sets the time when the maintenance period ends (Unix timestamp).
func WithMaintenanceActiveTill(activeTill int64) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.ActiveTill = StringInt64(activeTill)
	}
}

// WithMaintenanceDescription sets the description for the maintenance period.
func WithMaintenanceDescription(description string) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.Description = description
	}
}

// WithMaintenanceType sets the type of maintenance.
func WithMaintenanceType(maintenanceType MaintenanceType) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.MaintenanceType = maintenanceType
	}
}

// WithMaintenanceTimePeriods sets the time periods for the maintenance.
func WithMaintenanceTimePeriods(timePeriods []TimePeriod) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.TimePeriods = timePeriods
	}
}

// WithMaintenanceGroupIDs sets the host group IDs to be included in the maintenance.
func WithMaintenanceGroupIDs(groupIDs []string) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.GroupIDs = groupIDs
	}
}

// WithMaintenanceHostIDs sets the host IDs to be included in the maintenance.
func WithMaintenanceHostIDs(hostIDs []string) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.HostIDs = hostIDs
	}
}

// WithMaintenanceTags sets the problem tags for the maintenance.
func WithMaintenanceTags(tags []ProblemTag) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Params.Tags = tags
	}
}

// WithMaintenanceAuthToken sets the authentication token for the API request.
func WithMaintenanceAuthToken(token string) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.Auth = token
	}
}

// WithMaintenanceRequestID sets the ID for the API request.
func WithMaintenanceRequestID(id int) MaintenanceCreateOption {
	return func(mcr *MaintenanceCreateRequest) {
		mcr.ID = id
	}
}

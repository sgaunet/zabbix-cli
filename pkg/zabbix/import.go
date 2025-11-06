package zabbix

// API documentation: https://www.zabbix.com/documentation/7.2/en/manual/api/reference/configuration/importcompare

// const methodConfigurationImportCompare = "configuration.importcompare"

type importDiscoveryRules struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importGraphs struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importHostGroups struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importTemplateGroups struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importHosts struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importImages struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importMaps struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importMediaTypes struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importHTTPTests struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importItems struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importTemplateLinkage struct {
	CreateMissing bool `json:"createMissing"`
	DeleteMissing bool `json:"deleteMissing"`
}

type importTemplates struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
}

type importTemplateDashboards struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importTriggers struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

type importValueMaps struct {
	CreateMissing  bool `json:"createMissing"`
	UpdateExisting bool `json:"updateExisting"`
	DeleteMissing  bool `json:"deleteMissing"`
}

// rules defines import behavior for each object type in Zabbix configuration.
//
// Import Rules Structure (Zabbix API 7.2):
// Each object type has specific boolean flags that control the import behavior:
//
//   - CreateMissing (bool): Create objects that exist in import data but not in Zabbix
//   - UpdateExisting (bool): Update objects that exist in both import data and Zabbix
//   - DeleteMissing (bool): Delete objects that exist in Zabbix but not in import data
//
// Object Type Rules:
//
// Full Control (all three flags):
//   - discoveryRules: Low-level discovery rules
//   - graphs: Custom graphs for visualization
//   - httptests: Web monitoring scenarios
//   - items: Monitoring items and metrics
//   - templateDashboards: Dashboard definitions within templates
//   - triggers: Alert trigger definitions
//   - valueMaps: Value mapping definitions
//
// Create & Update Only (no DeleteMissing):
//   - host_groups: Host organization groups (Zabbix 6.2+)
//   - template_groups: Template organization groups (Zabbix 6.2+)
//   - hosts: Monitored host definitions
//   - images: Custom images for maps
//   - maps: Network topology maps
//   - mediaTypes: Notification channels
//   - templates: Reusable monitoring templates
//
// Special Behavior:
//   - templateLinkage: Has only CreateMissing and DeleteMissing (no UpdateExisting)
//     Controls linking of templates to hosts/templates
//
// Example Usage:
//
//	// Import with default rules (all flags set to true where applicable)
//	success, err := client.Import(ctx, yamlData)
//
//	// Custom rules can be configured by modifying the request before import
//	req := NewConfigurationImportRequest(yamlData)
//	// Modify req.Params.Rules as needed
//
// Note: Import operations validate the source data format and structure.
// Invalid data will result in API errors with detailed error messages.
type rules struct {
	DiscoveryRules     importDiscoveryRules     `json:"discoveryRules"`
	Graphs             importGraphs             `json:"graphs"`
	HostGroups         importHostGroups         `json:"host_groups"`
	TemplateGroups     importTemplateGroups     `json:"template_groups"`
	Hosts              importHosts              `json:"hosts"`
	HTTPTests          importHTTPTests          `json:"httptests"`
	Images             importImages             `json:"images"`
	Items              importItems              `json:"items"`
	Maps               importMaps               `json:"maps"`
	MediaTypes         importMediaTypes         `json:"mediaTypes"`
	TemplateLinkage    importTemplateLinkage    `json:"templateLinkage"`
	Templates          importTemplates          `json:"templates"`
	TemplateDashboards importTemplateDashboards `json:"templateDashboards"`
	Triggers           importTriggers           `json:"triggers"`
	ValueMaps          importValueMaps          `json:"valueMaps"`
}

type paramsImport struct {
	Format string `json:"format,omitempty"` // Optional: format of import data (yaml/json/xml)
	Rules  rules  `json:"rules,omitempty"`  // Optional: import rules (defaults applied by rulesAllTrue())
	Source string `json:"source"`           // Required: the import data content
}

type configurationImportCompareRequest struct {
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  paramsImport `json:"params"`
	Auth    string       `json:"auth"`
	ID      int          `json:"id"`
}

func rulesAllTrue() rules {
	return rules{
		DiscoveryRules:     importDiscoveryRules{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		Graphs:             importGraphs{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		HostGroups:         importHostGroups{CreateMissing: true, UpdateExisting: true},
		TemplateGroups:     importTemplateGroups{CreateMissing: true, UpdateExisting: true},
		Hosts:              importHosts{CreateMissing: true, UpdateExisting: true},
		HTTPTests:          importHTTPTests{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		Images:             importImages{CreateMissing: true, UpdateExisting: true},
		Items:              importItems{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		Maps:               importMaps{CreateMissing: true, UpdateExisting: true},
		MediaTypes:         importMediaTypes{CreateMissing: true, UpdateExisting: true},
		TemplateLinkage:    importTemplateLinkage{CreateMissing: true, DeleteMissing: true},
		Templates:          importTemplates{CreateMissing: true, UpdateExisting: true},
		TemplateDashboards: importTemplateDashboards{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		Triggers:           importTriggers{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		ValueMaps:          importValueMaps{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
	}
}

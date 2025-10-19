package zabbix

// API documentation: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/configuration/importcompare

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
	CreateMissing  bool `json:"createMissing"`
	DeleteExisting bool `json:"deleteMissing"`
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
	Format string `json:"format"`
	Rules  rules  `json:"rules"`
	Source string `json:"source"`
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
		TemplateLinkage:    importTemplateLinkage{CreateMissing: true, DeleteExisting: true},
		Templates:          importTemplates{CreateMissing: true, UpdateExisting: true},
		TemplateDashboards: importTemplateDashboards{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		Triggers:           importTriggers{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
		ValueMaps:          importValueMaps{CreateMissing: true, UpdateExisting: true, DeleteMissing: true},
	}
}

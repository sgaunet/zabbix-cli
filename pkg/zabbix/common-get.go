package zabbix

// More infos here: https://www.zabbix.com/documentation/6.0/en/manual/api/reference_commentary#common-get-method-parameters

type CommonGetParams struct {
	CountOutput   bool `json:"countOutput,omitempty"`
	Editable      bool `json:"editable,omitempty"`
	ExcludeSearch bool `json:"excludeSearch,omitempty"`
	// Filter 	object 	`json:"filter,omitempty"`
	Limit int `json:"limit,omitempty"`
	// Output 	query 	`json:"output,omitempty"`
	Preservekeys bool `json:"preservekeys,omitempty"`
	// Search 	object 	`json:"search,omitempty"`
	SearchByAny            bool     `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool     `json:"searchWildcardsEnabled,omitempty"`
	Sortorder              []string `json:"sortorder,omitempty"`
	StartSearch            bool     `json:"startSearch,omitempty"`
}

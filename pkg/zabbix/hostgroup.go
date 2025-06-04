package zabbix

// HostGroup represents the Zabbix host group API object.
// See: https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/object#host-group
type HostGroup struct {
	GroupID  string `json:"groupid,omitempty"`
	Name     string `json:"name"`
	Flags    string `json:"flags,omitempty"`    // Readonly. Origin of the host group: 0 - a plain host group; 4 - a discovered host group.
	Internal string `json:"internal,omitempty"` // Readonly. Whether the group is used internally by the Zabbix server. An internal group cannot be deleted. 0 - (default) not internal; 1 - internal.
	UUID     string `json:"uuid,omitempty"`     // Universal unique identifier, used for linking imported host groups to already existing ones.
}

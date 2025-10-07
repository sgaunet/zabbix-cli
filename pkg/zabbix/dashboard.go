package zabbix

// Dashboard represents the Zabbix dashboard API object.
// See: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/dashboard/object#dashboard
type Dashboard struct {
	DashboardID   string          `json:"dashboardid,omitempty"`
	Name          string          `json:"name"`
	UserID        string          `json:"userid,omitempty"`        // Owner of the dashboard.
	Private       string          `json:"private,omitempty"`       // Type of dashboard sharing. 0 - (default) public; 1 - private.
	DisplayPeriod string          `json:"display_period,omitempty"` // Display period in seconds. Default: 30.
	AutoStart     string          `json:"auto_start,omitempty"`    // Automatic slideshow. 0 - (default) disabled; 1 - enabled.
	Pages         []DashboardPage `json:"pages,omitempty"`         // Dashboard pages.
	Users         []DashboardUser `json:"users,omitempty"`         // Dashboard sharing with users.
	UserGroups    []DashboardUserGroup `json:"userGroups,omitempty"` // Dashboard sharing with user groups.
}

// DashboardPage represents a dashboard page.
// See: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/dashboard/object#dashboard-page
type DashboardPage struct {
	DashboardPageID string   `json:"dashboard_pageid,omitempty"`
	Name            string   `json:"name,omitempty"`
	DisplayPeriod   string   `json:"display_period,omitempty"` // Page display period in seconds. 0 - default host screen refresh rate. Default: 0.
	Widgets         []Widget `json:"widgets,omitempty"`        // Dashboard page widgets.
}

// Widget represents a dashboard widget.
// See: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/dashboard/object#widget
type Widget struct {
	WidgetID string        `json:"widgetid,omitempty"`
	Type     string        `json:"type"`                // Widget type.
	Name     string        `json:"name,omitempty"`      // Widget name.
	X        string        `json:"x,omitempty"`         // Widget horizontal position. 0 - leftmost. Default: 0.
	Y        string        `json:"y,omitempty"`         // Widget vertical position. 0 - topmost. Default: 0.
	Width    string        `json:"width,omitempty"`     // Widget width. Minimum: 1. Maximum: 24. Default: 1.
	Height   string        `json:"height,omitempty"`    // Widget height. Minimum: 2. Maximum: 32. Default: 2.
	View     string        `json:"view_mode,omitempty"` // Widget view mode. 0 - (default) normal; 1 - hidden header.
	Fields   []WidgetField `json:"fields,omitempty"`    // Widget configuration fields.
}

// WidgetField represents a single configuration field in a dashboard widget.
// See: https://www.zabbix.com/documentation/6.0/en/manual/api/reference/dashboard/object#widget-field
type WidgetField struct {
	Type  string      `json:"type"`            // Field type (determines how value is interpreted).
	Name  string      `json:"name"`            // Field name (e.g., "groupids", "severities", "tags").
	Value interface{} `json:"value,omitempty"` // Field value (can be string, int, or complex structure depending on field type).
}

// DashboardUser represents dashboard sharing with a user.
type DashboardUser struct {
	UserID     string `json:"userid"`
	Permission string `json:"permission"` // User permissions. 2 - (default) read-only; 3 - read-write.
}

// DashboardUserGroup represents dashboard sharing with a user group.
type DashboardUserGroup struct {
	UserGroupID string `json:"usrgrpid"`
	Permission  string `json:"permission"` // User group permissions. 2 - (default) read-only; 3 - read-write.
}

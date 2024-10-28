package zabbix

// Documentation of zabbix API: https://www.zabbix.com/documentation/6.0/en/manual/api/reference_commentary#common-get-method-parameters

// ZabbixFilter struct is used to filter the results from the Zabbix API (get method only)
// zabbix will return only those results that exactly match the given filter.
// Accepts an object, where the keys are property names, and the values are either a single value or an array of values to match against.
type ZabbixFilterGetMethod struct {
	filter map[string]interface{}
}

type zabbixFilterGetMethodOption func(*ZabbixFilterGetMethod)

// NewZabbixFilterGetMethod returns a new ZabbixFilterGetMethod
func NewZabbixFilterGetMethod(options ...zabbixFilterGetMethodOption) *ZabbixFilterGetMethod {
	z := &ZabbixFilterGetMethod{
		filter: make(map[string]interface{}),
	}
	for _, opt := range options {
		opt(z)
	}
	return z
}

// AddFilter adds a filter to the ZabbixFilterGetMethod
func AddFilter(key string, value interface{}) zabbixFilterGetMethodOption {
	return func(z *ZabbixFilterGetMethod) {
		z.filter[key] = value
	}
}

// GetFilter returns the filter
func (z *ZabbixFilterGetMethod) GetFilter() map[string]interface{} {
	return z.filter
}

// AddFilterName
func AddFilterName(name string) zabbixFilterGetMethodOption {
	return AddFilter("name", name)
}

// AddFilterNames
func AddFilterNames(names []string) zabbixFilterGetMethodOption {
	return AddFilter("name", names)
}

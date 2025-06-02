package zabbix

// Documentation of zabbix API: https://www.zabbix.com/documentation/6.0/en/manual/api/reference_commentary#common-get-method-parameters

// FilterGetMethod is used to filter the results from the Zabbix API (get method only).
// zabbix will return only those results that exactly match the given filter.
// Accepts an object, where the keys are property names, and the values are either a single value or an array of values to match against.
// FilterGetMethod holds filter parameters for get methods.
// FilterGetMethod holds filter parameters for get methods. (Linter: stutter is accepted for public API).
// FilterGetMethod holds filter parameters for get methods. (Linter: stutter is intentional for public API clarity).
type FilterGetMethod struct {
	filter map[string]interface{}
}

// FilterGetMethodOption is a function that modifies a FilterGetMethod.
// FilterGetMethodOption is a function that modifies a FilterGetMethod. (Linter: stutter is intentional for public API clarity).
type FilterGetMethodOption func(*FilterGetMethod)

// NewFilterGetMethod returns a new FilterGetMethod.
func NewFilterGetMethod(options ...FilterGetMethodOption) *FilterGetMethod {
	z := &FilterGetMethod{
		filter: make(map[string]interface{}),
	}
	for _, opt := range options {
		opt(z)
	}
	return z
}

// Filter adds a filter to the FilterGetMethod
// key is the property name
// value is the value to match against
// Be careful, the key must be a valid property name, no validation is done
// Filter returns a FilterGetMethodOption that sets a filter key and value.
func Filter(key string, value interface{}) FilterGetMethodOption {
	return func(z *FilterGetMethod) {
		z.filter[key] = value
	}
}

// GetFilter returns the filter.
func (z *FilterGetMethod) GetFilter() map[string]interface{} {
	return z.filter
}

// FilterByName returns a FilterGetMethodOption that filters by a single name.
func FilterByName(name string) FilterGetMethodOption {
	return Filter("name", name)
}

// FilterByNames returns a FilterGetMethodOption that filters by multiple names.
func FilterByNames(names []string) FilterGetMethodOption {
	return Filter("name", names)
}

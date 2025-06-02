package zabbix

// Severity represents the severity of an event.
type Severity int

// Severity values.
// Possible values:
// 0 - not classified;
// 1 - information;
// 2 - warning;
// 3 - average;
// 4 - high;
// 5 - disaster
const (
	NotClassified Severity = iota
	Information
	Warning
	Average
	High
	Disaster
)

// NewSeverity returns a Severity from an int value.
func NewSeverity(severity int) Severity {
	return Severity(severity)
}

// String returns the string representation of a severity.
func (s Severity) String() string {
	switch s {
	case NotClassified:
		return "Not classified"
	case Information:
		return "Information"
	case Warning:
		return "Warning"
	case Average:
		return "Average"
	case High:
		return "High"
	case Disaster:
		return "Disaster"
	default:
		return "Unknown"
	}
}

// GetSeverity returns the severity from an integer.
func GetSeverity(severity int) Severity {
	switch severity {
	case int(NotClassified):
		return NotClassified
	case int(Information):
		return Information
	case int(Warning):
		return Warning
	case int(Average):
		return Average
	case int(High):
		return High
	case int(Disaster):
		return Disaster
	default:
		return NotClassified
	}
}

// GetSeverityString returns the severity from a string.
func GetSeverityString(severity string) Severity {
	switch severity {
	case "Not classified":
		return NotClassified
	case "Information":
		return Information
	case "Warning":
		return Warning
	case "Average":
		return Average
	case "High":
		return High
	case "Disaster":
		return Disaster
	default:
		return NotClassified
	}
}

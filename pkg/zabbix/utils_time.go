package zabbix

import "time"

// FormatTimestamp converts a Unix timestamp to a human-readable date format.
func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

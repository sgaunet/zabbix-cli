package zabbix_test

import (
	"encoding/json"
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestTimePeriodTypeConstants(t *testing.T) {
	t.Parallel()

	// Verify all TimePeriodType constants match Zabbix API 7.2 spec
	tests := []struct {
		name     string
		typeVal  zabbix.TimePeriodType
		expected int
	}{
		{"OneTime", zabbix.TimePeriodTypeOneTime, 0},
		{"Daily", zabbix.TimePeriodTypeDaily, 1},
		{"Weekly", zabbix.TimePeriodTypeWeekly, 2},
		{"Monthly", zabbix.TimePeriodTypeMonthly, 3},
		{"MonthlyByWeekday", zabbix.TimePeriodTypeMonthlyByWeekday, 4},
		{"Yearly", zabbix.TimePeriodTypeYearly, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, int(tt.typeVal))
		})
	}
}

func TestTimePeriodOneTime(t *testing.T) {
	t.Parallel()

	// Type 0: One time only - requires start_date and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeOneTime,
		StartDate:      1717575000,
		Period:         3600,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(0), result["timeperiod_type"])
	require.Equal(t, float64(1717575000), result["start_date"])
	require.Equal(t, float64(3600), result["period"])
}

func TestTimePeriodDaily(t *testing.T) {
	t.Parallel()

	// Type 1: Daily - requires start_time and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeDaily,
		StartTime:      3600, // 01:00 AM
		Period:         7200, // 2 hours
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(1), result["timeperiod_type"])
	require.Equal(t, float64(3600), result["start_time"])
	require.Equal(t, float64(7200), result["period"])
}

func TestTimePeriodWeekly(t *testing.T) {
	t.Parallel()

	// Type 2: Weekly - requires dayofweek, start_time, and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeWeekly,
		DayOfWeek:      1, // Monday (0=Sunday, 1=Monday)
		StartTime:      3600,
		Period:         7200,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(2), result["timeperiod_type"])
	require.Equal(t, float64(1), result["dayofweek"])
	require.Equal(t, float64(3600), result["start_time"])
	require.Equal(t, float64(7200), result["period"])
}

func TestTimePeriodMonthly(t *testing.T) {
	t.Parallel()

	// Type 3: Monthly (by day of month) - requires day, start_time, and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeMonthly,
		Day:            15, // 15th of the month
		StartTime:      3600,
		Period:         7200,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(3), result["timeperiod_type"])
	require.Equal(t, float64(15), result["day"])
	require.Equal(t, float64(3600), result["start_time"])
	require.Equal(t, float64(7200), result["period"])
}

func TestTimePeriodMonthlyByWeekday(t *testing.T) {
	t.Parallel()

	// Type 4: Monthly (by day of week) - requires dayofweek, every, start_time, and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeMonthlyByWeekday,
		DayOfWeek:      1,    // Monday
		Every:          2,    // Second Monday of the month
		StartTime:      3600,
		Period:         7200,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(4), result["timeperiod_type"])
	require.Equal(t, float64(1), result["dayofweek"])
	require.Equal(t, float64(2), result["every"])
	require.Equal(t, float64(3600), result["start_time"])
	require.Equal(t, float64(7200), result["period"])
}

func TestTimePeriodYearly(t *testing.T) {
	t.Parallel()

	// Type 5: Yearly - requires month, day, start_time, and period
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeYearly,
		Month:          12, // December
		Day:            25, // 25th
		StartTime:      3600,  // 01:00 AM (not zero to avoid omitempty)
		Period:         86400, // 24 hours
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(5), result["timeperiod_type"])
	require.Equal(t, float64(12), result["month"])
	require.Equal(t, float64(25), result["day"])
	require.Equal(t, float64(3600), result["start_time"])
	require.Equal(t, float64(86400), result["period"])
}

func TestTimePeriodYearlyWithYear(t *testing.T) {
	t.Parallel()

	// Type 5: Yearly with optional year field
	tp := zabbix.TimePeriod{
		TimePeriodType: zabbix.TimePeriodTypeYearly,
		Month:          12,
		Day:            31,
		Year:           2025, // Optional year
		StartTime:      82800, // 11 PM
		Period:         7200,  // 2 hours
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Verify required fields are present including optional year
	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	require.NoError(t, err)

	require.Equal(t, float64(5), result["timeperiod_type"])
	require.Equal(t, float64(12), result["month"])
	require.Equal(t, float64(31), result["day"])
	require.Equal(t, float64(2025), result["year"])
	require.Equal(t, float64(82800), result["start_time"])
	require.Equal(t, float64(7200), result["period"])
}

func TestTimePeriodDayOfWeekValues(t *testing.T) {
	t.Parallel()

	// Test that DayOfWeek field is int (not array)
	tests := []struct {
		name      string
		dayOfWeek int
		expected  int
		omitted   bool // Fields with value 0 are omitted due to omitempty
	}{
		{"Sunday", 0, 0, true}, // 0 values are omitted with omitempty
		{"Monday", 1, 1, false},
		{"Tuesday", 2, 2, false},
		{"Wednesday", 3, 3, false},
		{"Thursday", 4, 4, false},
		{"Friday", 5, 5, false},
		{"Saturday", 6, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tp := zabbix.TimePeriod{
				TimePeriodType: zabbix.TimePeriodTypeWeekly,
				DayOfWeek:      tt.dayOfWeek,
				StartTime:      3600,
				Period:         7200,
			}

			jsonData, err := json.Marshal(tp)
			require.NoError(t, err)

			var result map[string]any
			err = json.Unmarshal(jsonData, &result)
			require.NoError(t, err)

			if tt.omitted {
				// Field with value 0 is omitted due to omitempty tag
				_, exists := result["dayofweek"]
				require.False(t, exists, "dayofweek should be omitted when value is 0")
			} else {
				require.Equal(t, float64(tt.expected), result["dayofweek"])
			}
		})
	}
}

func TestTimePeriodJSONMarshaling(t *testing.T) {
	t.Parallel()

	// Create a timeperiod with all fields
	tp := zabbix.TimePeriod{
		TimePeriodID:   "12345",
		TimePeriodType: zabbix.TimePeriodTypeMonthlyByWeekday,
		StartDate:      1717575000,
		Every:          2,
		Day:            15,
		DayOfWeek:      1,
		Month:          6,
		Year:           2024,
		StartTime:      3600,
		Period:         7200,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(tp)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled zabbix.TimePeriod
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Verify all fields are preserved
	require.Equal(t, tp.TimePeriodID, unmarshaled.TimePeriodID)
	require.Equal(t, tp.TimePeriodType, unmarshaled.TimePeriodType)
	require.Equal(t, tp.StartDate, unmarshaled.StartDate)
	require.Equal(t, tp.Every, unmarshaled.Every)
	require.Equal(t, tp.Day, unmarshaled.Day)
	require.Equal(t, tp.DayOfWeek, unmarshaled.DayOfWeek)
	require.Equal(t, tp.Month, unmarshaled.Month)
	require.Equal(t, tp.Year, unmarshaled.Year)
	require.Equal(t, tp.StartTime, unmarshaled.StartTime)
	require.Equal(t, tp.Period, unmarshaled.Period)
}

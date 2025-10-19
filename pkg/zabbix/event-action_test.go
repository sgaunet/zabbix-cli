package zabbix_test

import (
	"testing"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
	"github.com/stretchr/testify/require"
)

func TestNewEventAction(t *testing.T) {
	tests := []struct {
		name     string
		actions  []zabbix.EventAction
		expected int
	}{
		{
			name:     "Single action - CloseProblem",
			actions:  []zabbix.EventAction{zabbix.CloseProblem},
			expected: 1,
		},
		{
			name:     "Single action - Acknowledge",
			actions:  []zabbix.EventAction{zabbix.Acknowledge},
			expected: 2,
		},
		{
			name:     "Single action - AddMessage",
			actions:  []zabbix.EventAction{zabbix.AddMessage},
			expected: 4,
		},
		{
			name:     "Single action - ChangeSeverity",
			actions:  []zabbix.EventAction{zabbix.ChangeSeverity},
			expected: 8,
		},
		{
			name:     "Single action - Unacknowledge",
			actions:  []zabbix.EventAction{zabbix.Unacknowledge},
			expected: 16,
		},
		{
			name:     "Single action - Suppress",
			actions:  []zabbix.EventAction{zabbix.Suppress},
			expected: 32,
		},
		{
			name:     "Single action - Unsuppress",
			actions:  []zabbix.EventAction{zabbix.Unsuppress},
			expected: 64,
		},
		{
			name:     "Acknowledge and AddMessage (example from docs)",
			actions:  []zabbix.EventAction{zabbix.Acknowledge, zabbix.AddMessage},
			expected: 6, // 2 + 4
		},
		{
			name:     "AddMessage and ChangeSeverity (example from docs)",
			actions:  []zabbix.EventAction{zabbix.AddMessage, zabbix.ChangeSeverity},
			expected: 12, // 4 + 8
		},
		{
			name:     "All actions combined",
			actions:  []zabbix.EventAction{zabbix.CloseProblem, zabbix.Acknowledge, zabbix.AddMessage, zabbix.ChangeSeverity, zabbix.Unacknowledge, zabbix.Suppress, zabbix.Unsuppress},
			expected: 127, // 1 + 2 + 4 + 8 + 16 + 32 + 64
		},
		{
			name:     "CloseProblem and Acknowledge",
			actions:  []zabbix.EventAction{zabbix.CloseProblem, zabbix.Acknowledge},
			expected: 3, // 1 + 2
		},
		{
			name:     "No actions",
			actions:  []zabbix.EventAction{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zabbix.NewEventAction(tt.actions...)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRetrieveActions(t *testing.T) {
	tests := []struct {
		name     string
		action   zabbix.EventAction
		expected []zabbix.EventAction
	}{
		{
			name:     "CloseProblem only",
			action:   zabbix.CloseProblem,
			expected: []zabbix.EventAction{zabbix.CloseProblem},
		},
		{
			name:     "Acknowledge only",
			action:   zabbix.Acknowledge,
			expected: []zabbix.EventAction{zabbix.Acknowledge},
		},
		{
			name:     "AddMessage only",
			action:   zabbix.AddMessage,
			expected: []zabbix.EventAction{zabbix.AddMessage},
		},
		{
			name:     "ChangeSeverity only",
			action:   zabbix.ChangeSeverity,
			expected: []zabbix.EventAction{zabbix.ChangeSeverity},
		},
		{
			name:     "Unacknowledge only",
			action:   zabbix.Unacknowledge,
			expected: []zabbix.EventAction{zabbix.Unacknowledge},
		},
		{
			name:     "Suppress only",
			action:   zabbix.Suppress,
			expected: []zabbix.EventAction{zabbix.Suppress},
		},
		{
			name:     "Unsuppress only",
			action:   zabbix.Unsuppress,
			expected: []zabbix.EventAction{zabbix.Unsuppress},
		},
		{
			name:     "Acknowledge and AddMessage (6)",
			action:   zabbix.EventAction(6), // 2 + 4
			expected: []zabbix.EventAction{zabbix.Acknowledge, zabbix.AddMessage},
		},
		{
			name:     "AddMessage and ChangeSeverity (12)",
			action:   zabbix.EventAction(12), // 4 + 8
			expected: []zabbix.EventAction{zabbix.AddMessage, zabbix.ChangeSeverity},
		},
		{
			name:     "CloseProblem and Acknowledge (3)",
			action:   zabbix.EventAction(3), // 1 + 2
			expected: []zabbix.EventAction{zabbix.CloseProblem, zabbix.Acknowledge},
		},
		{
			name:     "All actions (127)",
			action:   zabbix.EventAction(127), // 1 + 2 + 4 + 8 + 16 + 32 + 64
			expected: []zabbix.EventAction{zabbix.CloseProblem, zabbix.Acknowledge, zabbix.AddMessage, zabbix.ChangeSeverity, zabbix.Unacknowledge, zabbix.Suppress, zabbix.Unsuppress},
		},
		{
			name:     "No actions (0)",
			action:   zabbix.EventAction(0),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zabbix.RetrieveActions(tt.action)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEventActionConstants(t *testing.T) {
	// Verify the bitmask values match the documentation
	require.Equal(t, 1, int(zabbix.CloseProblem))
	require.Equal(t, 2, int(zabbix.Acknowledge))
	require.Equal(t, 4, int(zabbix.AddMessage))
	require.Equal(t, 8, int(zabbix.ChangeSeverity))
	require.Equal(t, 16, int(zabbix.Unacknowledge))
	require.Equal(t, 32, int(zabbix.Suppress))
	require.Equal(t, 64, int(zabbix.Unsuppress))
}

func TestNewEventActionAndRetrieveRoundTrip(t *testing.T) {
	// Test that NewEventAction and RetrieveActions are inverse operations
	tests := []struct {
		name    string
		actions []zabbix.EventAction
	}{
		{
			name:    "Single action",
			actions: []zabbix.EventAction{zabbix.Acknowledge},
		},
		{
			name:    "Multiple actions",
			actions: []zabbix.EventAction{zabbix.Acknowledge, zabbix.AddMessage, zabbix.ChangeSeverity},
		},
		{
			name:    "All actions",
			actions: []zabbix.EventAction{zabbix.CloseProblem, zabbix.Acknowledge, zabbix.AddMessage, zabbix.ChangeSeverity, zabbix.Unacknowledge, zabbix.Suppress, zabbix.Unsuppress},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert actions to bitmask
			bitmask := zabbix.NewEventAction(tt.actions...)

			// Convert bitmask back to actions
			retrieved := zabbix.RetrieveActions(zabbix.EventAction(bitmask))

			// Should get the same actions back (order might differ, but content should match)
			require.ElementsMatch(t, tt.actions, retrieved)
		})
	}
}

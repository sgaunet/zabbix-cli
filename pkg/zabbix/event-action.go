package zabbix

// EventAction represents an event action type.
type EventAction int

// Event update action(s).
// Possible bitmap values are:
// 1 - close problem;
// 2 - acknowledge event;
// 4 - add message;
// 8 - change severity;
// 16 - unacknowledge event.
// This is a bitmask field; any sum of possible bitmap values is acceptable (for example, 6 for acknowledge event and add message).

const (
	// CloseProblem indicates the action to close a problem.
	CloseProblem   EventAction = 1
	// Acknowledge indicates the action to acknowledge an event.
	Acknowledge    EventAction = 2
	// AddMessage indicates the action to add a message to an event.
	AddMessage     EventAction = 4
	// ChangeSeverity indicates the action to change the severity of an event.
	ChangeSeverity EventAction = 8
	// Unacknowledge indicates the action to unacknowledge an event.
	Unacknowledge  EventAction = 16
)

// RettrieveActions returns the list of actions to perform on the event.
func RettrieveActions(action EventAction) []EventAction {
	var actions []EventAction
	if action&CloseProblem == CloseProblem {
		actions = append(actions, CloseProblem)
	}
	if action&Acknowledge == Acknowledge {
		actions = append(actions, Acknowledge)
	}
	if action&AddMessage == AddMessage {
		actions = append(actions, AddMessage)
	}
	if action&ChangeSeverity == ChangeSeverity {
		actions = append(actions, ChangeSeverity)
	}
	if action&Unacknowledge == Unacknowledge {
		actions = append(actions, Unacknowledge)
	}
	return actions
}

// NewEventAction returns the action to perform on the event.
func NewEventAction(actions ...EventAction) int {
	var action int
	for _, a := range actions {
		action |= int(a)
	}
	return action
}

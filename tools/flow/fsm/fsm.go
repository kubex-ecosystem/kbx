// Package fsm provides a simple finite state machine implementation.
//
// It defines generic State, Event, and Transition types, and a ThreadSafeFSM
// type that wraps an FSM with atomic operations for thread-safe state transitions.
//
// Example:
//
//	type State string
//	type Event string
//
//	const (
//	    Off State = "off"
//	    On  State = "on"
//	)
//
//	const (
//	    Toggle Event = "toggle"
//	)
//
//	transitions := []Transition{
//	    {From: Off, Event: Toggle, To: On},
//	    {From: On,  Event: Toggle, To: Off},
//	}
//
//	fsm := NewFSM(Off, transitions)
//	fsm.Trigger(Toggle) // transition from Off to On
package fsm

// State represents a state in the state machine.
type State string

// Event represents an event that can trigger a state transition.
type Event string

// Transition represents a transition between states.
type Transition struct {
	From  State `json:"from" yaml:"from" xml:"from" toml:"from" mapstructure:"from"`
	Event Event `json:"event" yaml:"event" xml:"event" toml:"event" mapstructure:"event"`
	To    State `json:"to" yaml:"to" xml:"to" toml:"to" mapstructure:"to"`
}

// FSM represents a finite state machine.
type FSM struct {
	current State                     `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
	table   map[State]map[Event]State `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

// NewFSM creates a new FSM with the given initial state and transitions.
func NewFSM(initial State, transitions []Transition) *FSM {
	fsm := &FSM{
		current: initial,
		table:   make(map[State]map[Event]State),
	}

	for _, t := range transitions {
		if fsm.table[t.From] == nil {
			fsm.table[t.From] = make(map[Event]State)
		}
		fsm.table[t.From][t.Event] = t.To
	}

	return fsm
}

// Current returns the current state of the FSM.
func (f *FSM) Current() State {
	return f.current
}

// Trigger transitions the FSM to a new state based on the given event.
func (f *FSM) Trigger(event Event) bool {
	if next, ok := f.table[f.current][event]; ok {
		f.current = next
		return true
	}
	return false
}

// Can checks if the given event can trigger a transition from the current state.
func (f *FSM) Can(event Event) bool {
	_, ok := f.table[f.current][event]
	return ok
}

// Reset resets the FSM to the given state.
func (f *FSM) Reset(state State) {
	f.current = state
}

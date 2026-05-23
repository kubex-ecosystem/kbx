package fsm

import "github.com/kubex-ecosystem/kbx/tools/flow/control"

// State represents an atomic state based on bitflags.
type State uint32

// Event represents an event that can trigger a state transition.
type Event uint32

// Transition represents an atomic transition between states.
type Transition struct {
	From  State
	Event Event
	To    State
}

// FSM is an atomic finite state machine.
type FSM struct {
	current control.FlagReg32[State]
	table   map[State]map[Event]State
}

// NewFSM creates a new atomic FSM with the given initial state and transitions.
func NewFSM(initial State, transitions []Transition) *FSM {
	f := &FSM{table: make(map[State]map[Event]State)}
	f.current.Store(initial)
	for _, t := range transitions {
		if f.table[t.From] == nil {
			f.table[t.From] = make(map[Event]State)
		}
		f.table[t.From][t.Event] = t.To
	}
	return f
}

// Current returns the current state.
func (f *FSM) Current() State { return f.current.Load() }

// Trigger transitions the FSM atomically. Returns false if the event is invalid for the current state.
func (f *FSM) Trigger(event Event) bool {
	for {
		curr := f.current.Load()
		next, ok := f.table[curr][event]
		if !ok {
			return false
		}
		if f.current.CompareAndSwap(curr, next) {
			return true
		}
	}
}

// Can reports whether the given event can trigger a transition from the current state.
func (f *FSM) Can(event Event) bool {
	_, ok := f.table[f.current.Load()][event]
	return ok
}

// Reset forces the FSM to the given state.
func (f *FSM) Reset(state State) { f.current.Store(state) }

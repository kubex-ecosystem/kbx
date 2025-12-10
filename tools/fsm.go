package tools

type State string

type Event string

type Transition struct {
	From  State
	Event Event
	To    State
}

type FSM struct {
	current State
	table   map[State]map[Event]State
}

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

func (f *FSM) Current() State {
	return f.current
}

func (f *FSM) Trigger(event Event) bool {
	if next, ok := f.table[f.current][event]; ok {
		f.current = next
		return true
	}
	return false
}

func (f *FSM) Can(event Event) bool {
	_, ok := f.table[f.current][event]
	return ok
}

func (f *FSM) Reset(state State) {
	f.current = state
}

package flow

import "github.com/kubex-ecosystem/kbx/tools/flow/fsm"

// FlowState tracks the lifecycle of a flow execution.
type FlowState struct {
	started   bool
	completed bool
	failed    bool
}

func (s *FlowState) Start() error {
	s.started = true
	return nil
}

func (s *FlowState) Complete() {
	s.completed = true
}

func (s *FlowState) Fail() {
	s.failed = true
}

// Flow bundles an FSM and a lifecycle state for a single orchestration unit.
type Flow struct {
	FSM   *fsm.FSM
	State FlowState
}

// NewFlow creates a new Flow instance.
func NewFlow() *Flow {
	return &Flow{}
}

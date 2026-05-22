package flow

import (
	"context"
	"sync"

	"github.com/kubex-ecosystem/kbx/tools/flow/control"
	"github.com/kubex-ecosystem/kbx/tools/flow/events"
	"github.com/kubex-ecosystem/kbx/tools/flow/fsm"
	"github.com/kubex-ecosystem/kbx/tools/flow/queue"
)

// Flow is the high-level orchestrator.
type Flow struct {
	mu    sync.RWMutex
	State control.JobState
	Queue *queue.Queue
	FSM   *fsm.FSM
}

func NewFlow() *Flow {
	return &Flow{
		Queue: queue.NewQueue(100),
	}
}

func (f *Flow) Start(workers int) {
	f.Queue.Start(workers)
}

func (f *Flow) Submit(ctx context.Context, task queue.Task) error {
	if err := f.State.Start(); err != nil {
		return err
	}
	f.Queue.Enqueue(task)
	return nil
}

// Alias for tests/compatibility
type LogHooks = events.LogHooks
type HookFunc[T any] = events.HookFunc[T]

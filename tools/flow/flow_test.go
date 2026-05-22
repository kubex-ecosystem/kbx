package flow_test

import (
	"context"
	"github.com/kubex-ecosystem/kbx/tools/flow"
	"github.com/kubex-ecosystem/kbx/tools/flow/control"
	"testing"
)

func TestFlowAtomicStart(t *testing.T) {
	f := flow.NewFlow()
	ctx := context.Background()

	// Task dummy
	task := func(ctx context.Context) error {
		return nil
	}

	// First submission should work
	err := f.Submit(ctx, task)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	if !f.State.Has(control.JobRunning) {
		t.Error("Expected state to have JobRunning flag")
	}

	// Marking as complete
	err = f.State.Complete()
	if err != nil {
		t.Fatalf("Expected nil error on complete, got %v", err)
	}

	if !f.State.IsTerminal() {
		t.Error("Expected state to be terminal")
	}

	// Submission after terminal should fail
	err = f.Submit(ctx, task)
	if err == nil {
		t.Error("Expected error when submitting to terminal flow, got nil")
	}
}

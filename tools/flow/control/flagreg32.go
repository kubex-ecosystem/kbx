package control

import (
	"errors"
	"sync/atomic"
)

// FlagReg32 is a generic atomic wrapper for uint32-based flag registers.
type FlagReg32[T ~uint32] struct {
	val atomic.Uint32
}

func (r *FlagReg32[T]) Store(v T) {
	r.val.Store(uint32(v))
}

func (r *FlagReg32[T]) Load() T {
	return T(r.val.Load())
}

func (r *FlagReg32[T]) CompareAndSwap(old, new T) bool {
	return r.val.CompareAndSwap(uint32(old), uint32(new))
}

// StepFlag represents a pipeline processing stage.
type StepFlag uint32

const (
	StepIdle      StepFlag = iota
	StepValidate  StepFlag = 1 << iota
	StepPreHooks  StepFlag = 1 << iota
	StepFormat    StepFlag = 1 << iota
	StepPostHooks StepFlag = 1 << iota
	StepWrite     StepFlag = 1 << iota
	StepDone      StepFlag = 1 << iota
	StepFailed    StepFlag = 1 << iota
)

// ErrTerminal is returned when processing is attempted after the pipeline has terminated.
var ErrTerminal = errors.New("manager: pipeline has terminated")

// ManagerControl tracks the atomic processing state of a pipeline.
type ManagerControl struct {
	Stage FlagReg32[StepFlag]
}

// IsTerminal reports whether the pipeline is in a terminal state (done or failed).
func (c *ManagerControl) IsTerminal() bool {
	s := c.Stage.Load()
	return s == StepDone || s == StepFailed
}

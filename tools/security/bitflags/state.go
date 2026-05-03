// Package bitflags provides atomic flag registers and utilities for managing
// sets of bit flags on top of basic integer types.
//
// The package defines generic flag types `FlagReg32[T ~uint32]` and
// `FlagReg64[T ~uint64]`, along with helper utilities. The type parameter `T`
// must be an unsigned integer type.
//
// Operations like Set, Clear, Toggle, and Mask are implemented using atomic
// compare-and-swap (CAS) loops, ensuring thread safety without requiring
// external mutexes.
//
// Conditional helpers `SetIf` and `ClearIf` allow modifying flags only when
// certain other flags are set or clear respectively, with atomic guarantees.
//
// Pretty-print utilities provide formatted string representations of flag
// combinations, mapping flag values to symbolic names.
//
// Example:
//
//	type Perms uint32
//
//	const (
//	    Read  Perms = 1 << iota
//	    Write
//	    Exec
//	)
//
//	var p FlagReg32[Perms]
//	p.Set(Read | Write)
//	if p.Has(Read) { fmt.Println("readable") }
package bitflags

import "errors"

// Example: Job state flags (combináveis) -------------------------------------

type JobFlag uint32

const (
	JobPending JobFlag = 1 << iota
	JobRunning
	JobCancelRequested
	JobRetrying
	JobCompleted
	JobFailed
	JobTimedOut
)

const (
	terminalMask JobFlag = JobCompleted | JobFailed | JobTimedOut
)

var (
	ErrTerminal = errors.New("job is in a terminal state")
)

type JobState struct{ r FlagReg32[JobFlag] }

func (s *JobState) Load() JobFlag { return s.r.Load() }

// Start only from Pending; sets Running.
func (s *JobState) Start() error {
	ok := s.r.SetIf(terminalMask|JobRunning|JobCompleted|JobFailed|JobTimedOut, JobRunning)
	if !ok {
		return ErrTerminal
	}
	return nil
}

func (s *JobState) RequestCancel() { s.r.Set(JobCancelRequested) }
func (s *JobState) Retry() error {
	// can retry if not terminal; set Retrying and clear Running
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobRetrying) &^ JobRunning
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) Complete() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobCompleted) &^ (JobRunning | JobRetrying | JobCancelRequested)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) Fail() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobFailed) &^ (JobRunning | JobRetrying)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) Timeout() error {
	for {
		old := s.r.Load()
		if old&terminalMask != 0 {
			return ErrTerminal
		}
		newV := (old | JobTimedOut) &^ (JobRunning | JobRetrying)
		if s.r.CompareAndSwap(old, newV) {
			return nil
		}
	}
}

func (s *JobState) IsTerminal() bool { return s.r.Any(terminalMask) }

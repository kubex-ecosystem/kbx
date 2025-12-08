package kbx

import "fmt"

type IBitstate[T ~uint64] interface {
	Set(flag T)
	Clear(flag T)
	Has(flag T) bool
	WaitFor(flag T)
}

type Bitstate[T ~uint64] struct {
	state T
}

func NewBitstate[T ~uint64]() *Bitstate[T] {
	return &Bitstate[T]{state: 0}
}

func (b *Bitstate[T]) Set(flag T) {
	b.state |= flag
}

func (b *Bitstate[T]) Clear(flag T) {
	b.state &^= flag
}

func (b *Bitstate[T]) Has(flag T) bool {
	return b.state&flag != 0
}

func (b *Bitstate[T]) WaitFor(flag T) {
	for !b.Has(flag) {
		// Busy-waiting; in a real implementation, consider using sync.Cond or channels
	}
}

func (b *Bitstate[T]) State() T {
	return b.state
}

func (b *Bitstate[T]) Reset() {
	b.state = 0
}

func (b *Bitstate[T]) Toggle(flag T) {
	b.state ^= flag
}

func (b *Bitstate[T]) IsEmpty() bool {
	return b.state == 0
}

func (b *Bitstate[T]) SetAll(flags T) {
	b.state = flags
}

func (b *Bitstate[T]) ClearAll() {
	b.state = 0
}

func (b *Bitstate[T]) Copy() *Bitstate[T] {
	return &Bitstate[T]{state: b.state}
}

func (b *Bitstate[T]) String() string {
	return fmt.Sprintf("%b", b.state)
}

func (b *Bitstate[T]) Equals(other *Bitstate[T]) bool {
	return b.state == other.state
}

func (b *Bitstate[T]) And(other *Bitstate[T]) T {
	return b.state & other.state
}

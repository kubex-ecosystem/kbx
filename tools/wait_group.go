package tools

import "sync"

// GoGroup is a simple wrapper around sync.WaitGroup to manage goroutines.
type GoGroup struct{ wg sync.WaitGroup }

type GoFunc[T any] func() T

// Go starts a new goroutine and tracks it in the WaitGroup.
func (g *GoGroup) Go(fn GoFunc[any]) {
	if g != nil {
		g.wg.Go(func() {
			fn()
		})
		g.Wait()
	}
}
func (g *GoGroup) Wait() {
	if g != nil {
		g.wg.Wait()
	} else {
		return
	}
}

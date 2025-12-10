package tools

import "sync"

type GoGroup struct {
	wg sync.WaitGroup
}

func (g *GoGroup) Go(fn func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fn()
	}()
}

func (g *GoGroup) Wait() {
	g.wg.Wait()
}

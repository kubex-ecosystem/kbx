package tools

import "context"

type Task func(context.Context) error

type Queue struct {
	tasks chan Task
}

func NewQueue(size int) *Queue {
	return &Queue{tasks: make(chan Task, size)}
}

func (q *Queue) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			for t := range q.tasks {
				t(context.Background())
			}
		}()
	}
}

func (q *Queue) Enqueue(t Task) {
	q.tasks <- t
}

func (q *Queue) Stop() {
	close(q.tasks)
}

func (q *Queue) Size() int {
	return len(q.tasks)
}

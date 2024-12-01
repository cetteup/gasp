package task

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

type Task func(ctx context.Context) error

type AsyncRunner struct {
	tasks []Task
	mu    sync.Mutex
	err   error
}

func (r *AsyncRunner) Append(tasks ...Task) {
	for _, task := range tasks {
		if task != nil {
			r.tasks = append(r.tasks, task)
		}
	}
}

func (r *AsyncRunner) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, task := range r.tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := task(ctx)
			if err != nil {
				r.mu.Lock()
				r.err = multierr.Append(r.err, err)
				r.mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return r.err
}

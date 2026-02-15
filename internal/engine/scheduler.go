package engine

import (
	"context"
	"sync"

	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Scheduler struct {
	ctx         context.Context
	concurrency int

	tasks   chan tasks.Task
	results chan domain.Result

	wg sync.WaitGroup
}

func NewScheduler(ctx context.Context, concurrency int, results chan domain.Result) *Scheduler {
	return &Scheduler{
		ctx:         ctx,
		concurrency: concurrency,
		tasks:       make(chan tasks.Task),
		results:     results,
	}
}

func (s *Scheduler) Run(ts []tasks.Task) {

	s.createWorkers()

	go func() {
		s.assignTasks(ts)
		close(s.tasks)
	}()

	s.wg.Wait()
	close(s.results)
}

func (s *Scheduler) assignTasks(ts []tasks.Task) {
	for _, t := range ts {
		select {
		case s.tasks <- t:
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) createWorkers() {
	for i := 0; i < s.concurrency; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Scheduler) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case task, ok := <-s.tasks:

			if !ok {
				return
			}

			_ = task.Execute(s.ctx, s.results)
		}
	}
}

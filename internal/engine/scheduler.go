package engine

import (
	"context"
	"sync"

	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Scheduler struct {
	ctx                  context.Context
	concurrencyThreshold int

	tasks   chan tasks.Task
	results chan domain.Result

	wg sync.WaitGroup
}

func NewScheduler(ctx context.Context, results chan domain.Result, treshhold int) *Scheduler {

	return &Scheduler{
		ctx:                  ctx,
		concurrencyThreshold: treshhold,
		tasks:                make(chan tasks.Task),
		results:              results,
	}
}

func (s *Scheduler) Run(ts []tasks.Task) {

	numWorkers := s.concurrencyThreshold
	numTasks := len(ts)

	if numTasks < s.concurrencyThreshold {

		numWorkers = numTasks
	}

	s.createWorkers(numWorkers)

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

func (s *Scheduler) createWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
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

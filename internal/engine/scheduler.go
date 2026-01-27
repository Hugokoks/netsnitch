package engine

import (
	"context"
	"sync"

	"netsnitch/internal/scan"
	"netsnitch/internal/scanner"
	"netsnitch/internal/target"
)

type Scheduler struct {
	concurrency int
	targets     chan target.Target
	results     chan scan.Result
	scanner     scanner.Scanner
	wg          sync.WaitGroup
	ctx         context.Context
}

func NewScheduler(ctx context.Context, concurrency int, sc scanner.Scanner) *Scheduler {
	return &Scheduler{
		ctx:         ctx,
		concurrency: concurrency,
		scanner:     sc,
		results:     make(chan scan.Result, 100),
		targets:     make(chan target.Target, 100),
	}
}

func (s *Scheduler) RunScan(targets []target.Target) {
	s.start()

	s.add(targets)

	s.end()
}

// //create workers gourutines
func (s *Scheduler) start() {
	for i := 0; i < s.concurrency; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// //Pass task to workers
func (s *Scheduler) add(targets []target.Target) {

	for _, t := range targets {
		select {
		case s.targets <- t:
		case <-s.ctx.Done():
			return
		}
	}
}

// Results returns a read-only channel with scan results.
// Exposing the channel as receive-only prevents external consumers
// from sending to or closing the channel, which is managed internally
// by the Scheduler.
func (s *Scheduler) Results() <-chan scan.Result {

	return s.results
}

// ///End's worker gourutines and close task chan
func (s *Scheduler) end() {
	close(s.targets)
	s.wg.Wait()
	close(s.results)
}

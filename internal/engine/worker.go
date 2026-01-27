package engine

///// 1. waiting for task from chann s.tasks if there is no task wait
///// 2. do task
///// 3. close wait group
////  4. all wait groups done close task conn -> end goroutine

func (s *Scheduler) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case target, ok := <-s.targets:
			if !ok {
				return
			}

			res := s.scanner.Scan(s.ctx, target)

			select {
			case s.results <- res:
			case <-s.ctx.Done():
				return
			}
		}
	}
}

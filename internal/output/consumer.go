package output

import (
	"context"
	"netsnitch/internal/domain"
	"sync"
)

type Consumer struct {
	wg      sync.WaitGroup
	ctx     context.Context
	results chan domain.Result
}

func NewConsumer(ctx context.Context, results chan domain.Result) Consumer {
	return Consumer{
		ctx:     ctx,
		results: results,
	}
}

func (c *Consumer) Start() {

	c.wg.Add(1)
	go c.Consume()

}
func (c *Consumer) Wait() {

	c.wg.Wait()
}

func (c *Consumer) Consume() {

	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return

		case res, ok := <-c.results:
			if !ok {
				return
			}

			out(res)

		}
	}
}

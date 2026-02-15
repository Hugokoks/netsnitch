package engine

import (
	"context"
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
	"netsnitch/internal/tasks"
)

func Run(ctx context.Context, ts []tasks.Task, concurrency int) error {
	fmt.Println("[engine] starting scan")

	results := make(chan domain.Result, 100)

	scheduler := NewScheduler(ctx, concurrency, results)
	consumer := output.NewConsumer(ctx, results)

	go consumer.Start()
	scheduler.Run(ts)

	consumer.Wait()
	fmt.Println("[engine] scan finished")
	return nil
}

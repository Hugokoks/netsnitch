package engine

import (
	"context"
	"fmt"
	"netsnitch/internal/output"
	"netsnitch/internal/tasks"
)

func Run(ctx context.Context, ts []tasks.Task, concurrency int) error {
	fmt.Println("[engine] starting scan")

	scheduler := NewScheduler(ctx, concurrency)

	go output.ConsumeResults(ctx, scheduler.results)

	scheduler.Run(ts)

	fmt.Println("[engine] scan finished")
	return nil
}

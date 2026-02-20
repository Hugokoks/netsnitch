package engine

import (
	"context"
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/output"
	"netsnitch/internal/tasks"
	"os"
	"strconv"
)

func Run(ctx context.Context, ts []tasks.Task) error {
	fmt.Println("[engine] starting scan")

	concurrencyStr := os.Getenv("CONCURRENCY_THRESHOLD")
	threshold, err := strconv.Atoi(concurrencyStr)
	if err != nil {
		threshold = 200
	}

	results := make(chan domain.Result, 100)

	scheduler := NewScheduler(ctx, results, threshold)

	consumer := output.NewConsumer(ctx, results)

	go consumer.Start()
	scheduler.Run(ts)

	consumer.Wait()
	fmt.Println("[engine] scan finished")
	return nil
}

package tasks

import (
	"context"
	"netsnitch/internal/domain"
)

type Task interface {
	Execute(ctx context.Context, out chan<- domain.Result) error
}


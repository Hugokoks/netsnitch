package tcp

import (
	"context"
	"net"
	"netsnitch/internal/domain"
	"time"
)

type Task struct {
	timeout time.Duration
	ip      net.IP
	port    int
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error {

	/////TO DO: try to remove this select blog and see how context is behaving after CTRL + C
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	res := scanTarget(ctx, t.ip, t.port, t.timeout)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- res:
		return nil
	}
}

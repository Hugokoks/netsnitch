package tcp

import (
	"context"
	"net"
	"time"

	"netsnitch/internal/domain"
)

type Task struct {
	ip    net.IP
	port   int
	timeout time.Duration
}

func NewTask(ip net.IP, port int, timeout time.Duration) *Task {
	return &Task{
		ip:      ip,
		port:    port,
		timeout: timeout,
	}
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	res := scanTarget(ctx, t.ip, t.port, t.timeout)
	out <- res
	return nil
}
package tcp

import (
	"context"
	"net"
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/tcp/tcp_full"
	"netsnitch/internal/scans/tcp/tcp_stealth"
	"time"
)

type baseTask struct {
	ip       net.IP
	port     int
	timeout  time.Duration
	render   domain.RenderType
	openOnly bool
}

func (t *baseTask) sendResult(ctx context.Context, res domain.Result, out chan<- domain.Result) error {
	if t.openOnly && !res.Open {
		return nil
	}
	res.RenderType = t.render
	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- res:
		return nil
	}
}

type StealthTask struct {
	baseTask
	mgr *tcp_stealth.Manager
}

func (t *StealthTask) Execute(ctx context.Context, out chan<- domain.Result) error {
	res := t.mgr.Scan(ctx, t.ip, t.port, t.timeout)
	return t.sendResult(ctx, res, out)
}

type FullTask struct {
	baseTask
}

func (t *FullTask) Execute(ctx context.Context, out chan<- domain.Result) error {
	res := tcp_full.Scan(ctx, t.ip, t.port, t.timeout)
	return t.sendResult(ctx, res, out)
}

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
	mode    domain.ScanMode
	mgr     *StealthManager
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error {

	/////TODO: Make switch based on scan mode - full, stealth, ...

	var res domain.Result
	switch t.mode {

	case domain.STEALTH:
		res = t.mgr.Scan(t.ip, t.port, t.timeout)
	case domain.FULL:
		res = fullScan(ctx, t.ip, t.port, t.timeout)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case out <- res:
		return nil
	}
}

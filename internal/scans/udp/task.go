package udp

import (
	"context"
	"net"
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/udp/udp_scan"
	"time"
)

type UDPTask struct {
	ip       net.IP
	port     int
	timeout  time.Duration
	render   domain.RenderType
	openOnly bool
	mgr      *udp_scan.Manager
}

func (t *UDPTask) Execute(ctx context.Context, out chan<- domain.Result) error {

	res := t.mgr.Scan(ctx, t.ip, t.port, t.timeout)

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

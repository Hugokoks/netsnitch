package udp

import (
	"context"
	"net"
	"netsnitch/internal/domain"
	"time"
)

type UDPTask struct {
	ip       net.IP
	port     int
	timeout  time.Duration
	render   domain.RenderType
	openOnly bool
}

func (t *UDPTask) Execute(ctx context.Context, out chan<- domain.Result) error {

	res := domain.Result{}

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

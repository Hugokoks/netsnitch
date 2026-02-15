package arp_active

import (
	"context"
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/arp/discovery"

	"time"
)

type Task struct {
	timeout time.Duration
	scope   domain.Scope
	mode    domain.ScanMode
	render  domain.RenderType
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error {

	ips, err := domain.ResolveScope(t.scope)
	if err != nil {
		return err
	}

	arp := discovery.NewARP(t.timeout)

	res, err := arp.Discover(ctx, ips, t.mode)
	if err != nil {
		return err
	}

	for _, r := range res {
		select {

		case <-ctx.Done():
			return ctx.Err()

		case out <- domain.Result{
			Protocol: domain.ARP,
			IP:       r.IP,
			MAC:      r.MAC,
			Alive:    true,
		}:

		}
	}

	return nil
}

package arp_active

import (
	"context"
	"netsnitch/internal/discovery"
	"netsnitch/internal/domain"
	"time"
)

type Task struct {
	timeout time.Duration
	scope   domain.Scope
	Mode    domain.ScanMode
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error {

	ips, err := domain.ResolveScope(t.scope)
	if err != nil {
		return err
	}

	arp := discovery.NewARP(t.timeout)

	res, err := arp.Discover(ctx, ips, t.Mode)
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

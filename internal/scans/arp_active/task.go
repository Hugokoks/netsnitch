package arp_active

import (
	"context"
	"netsnitch/internal/discovery"
	"netsnitch/internal/domain"
	"time"
)

type Task struct{
	timeout time.Duration
	cidr string
}

func (t *Task) Execute(ctx context.Context, out chan<- domain.Result) error{


	arp := discovery.NewARP(t.timeout)

	res, err := arp.Discover(ctx, t.cidr, domain.ARP_ACTIVE)
	
	if err != nil {
		return err
	}

	/////Send results to chan
	for _,r := range res{
		select{

			case <-ctx.Done():
				return ctx.Err()

			case out <- domain.Result{
				Protocol: domain.ARP_ACTIVE,
				IP: r.IP,
				MAC:r.MAC,
				Alive: true,
			}:
		
		}
	
	}

	return nil
}
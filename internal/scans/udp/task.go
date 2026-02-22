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

	result := domain.Result{
		Protocol: domain.UDP,
		IP:       t.ip,
		Port:     t.port,
	}

	for _, probe := range udpProbes {

		if len(probe.Ports) > 0 && !contains(probe.Ports, t.port) {
			continue
		}

		state := t.mgr.Scan(ctx, t.ip, t.port, t.timeout, probe.Payload)

		switch state {
		case udp_scan.UDPOpen:
			result.Open = true
			goto SEND
		case udp_scan.UDPClosed:
			result.Open = false
			goto SEND
		case udp_scan.UDPOpenOrFiltered:

		}
	}
	result.Service = "open|filtered"

SEND:
	if t.openOnly && !result.Open {
		return nil
	}
	result.RenderType = t.render
	out <- result
	return nil
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

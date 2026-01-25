package discovery

import (
	"context"
	"net"
	"time"
)

type ARPDiscoverer struct {
	Timeout time.Duration
}

func NewARP(timeout time.Duration) *ARPDiscoverer {
	return &ARPDiscoverer{Timeout: timeout}
}

func (a *ARPDiscoverer) Discover(ctx context.Context, cidr string) ([]net.IP, error) {
	// 1. Parse CIDR
	ips, err := parseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// 2. Pick interface + source IP
	iface, srcIP, err := pickInterface(ips)
	if err != nil {
		return nil, err
	}

	// 3. Open ARP handle
	handle, err := openARPHandle(iface)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	// 4. Send ARP requests
	for _, ip := range ips {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			_ = sendARPRequest(handle, iface, srcIP, ip)
		}
	}

	// 5. Collect replies
	alive := make(map[string]net.IP)
	timeout := time.After(a.Timeout)

	for {
		select {
		case <-ctx.Done():
			return mapToSlice(alive), nil

		case <-timeout:
			return mapToSlice(alive), nil

		default:
			if ip := readARPReply(handle); ip != nil {
				alive[ip.String()] = ip
			}
		}
	}
}

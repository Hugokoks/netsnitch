package discovery

import (
	"context"
	"net"
)

// Discoverer finds alive hosts in a network.
// It returns a slice of IP addresses that responded.

type Discoverer interface {
	Discover(ctx context.Context, cidr string) ([]net.IP, error)
}

package domain

import (
	"fmt"
	"net"
)

type ScopeType int

const (
	ScopeCIDR ScopeType = iota
	ScopeIPs
)

type Scope struct {
	Type ScopeType

	CIDR string
	IPs  []net.IP
}

func ResolveScope(s Scope) ([]net.IP, error) {
	switch s.Type {
	case ScopeCIDR:
		return ParseCIDR(s.CIDR)
	case ScopeIPs:
		return s.IPs, nil
	default:
		return nil, fmt.Errorf("unknown scope type")
	}
}

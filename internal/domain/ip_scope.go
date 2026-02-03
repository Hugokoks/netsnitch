package domain

import (
	"fmt"
	"net"
	"strings"
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

func ParseScope(token string) (Scope, error) {

	token = strings.TrimSpace(token)

	// More IP's
	if strings.Contains(token, ",") {
		parts := strings.Split(token, ",")
		var ips []net.IP

		for _, p := range parts {
			ip := net.ParseIP(strings.TrimSpace(p))
			if ip == nil {
				return Scope{}, fmt.Errorf("invalid IP: %s", p)
			}
			ips = append(ips, ip)
		}

		return Scope{
			Type: ScopeIPs,
			IPs:  ips,
		}, nil
	}

	// CIDR
	if strings.Contains(token, "/") {
		if _, _, err := net.ParseCIDR(token); err != nil {
			return Scope{}, fmt.Errorf("invalid CIDR: %s", token)
		}

		return Scope{
			Type: ScopeCIDR,
			CIDR: token,
		}, nil
	}

	// single IP
	ip := net.ParseIP(token)
	if ip == nil {
		return Scope{}, fmt.Errorf("invalid IP or CIDR: %s", token)
	}

	return Scope{
		Type: ScopeIPs,
		IPs:  []net.IP{ip},
	}, nil
}

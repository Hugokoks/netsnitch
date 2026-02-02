package input

import (
	"fmt"
	"net"
	"strings"

	"netsnitch/internal/domain"
)

func ParseScope(token string) (domain.Scope, error) {

	token = strings.TrimSpace(token)

	// More IP's
	if strings.Contains(token, ",") {
		parts := strings.Split(token, ",")
		var ips []net.IP

		for _, p := range parts {
			ip := net.ParseIP(strings.TrimSpace(p))
			if ip == nil {
				return domain.Scope{}, fmt.Errorf("invalid IP: %s", p)
			}
			ips = append(ips, ip)
		}

		return domain.Scope{
			Type: domain.ScopeIPs,
			IPs:  ips,
		}, nil
	}

	// CIDR
	if strings.Contains(token, "/") {
		if _, _, err := net.ParseCIDR(token); err != nil {
			return domain.Scope{}, fmt.Errorf("invalid CIDR: %s", token)
		}

		return domain.Scope{
			Type: domain.ScopeCIDR,
			CIDR: token,
		}, nil
	}

	// single IP
	ip := net.ParseIP(token)
	if ip == nil {
		return domain.Scope{}, fmt.Errorf("invalid IP or CIDR: %s", token)
	}

	return domain.Scope{
		Type: domain.ScopeIPs,
		IPs:  []net.IP{ip},
	}, nil
}

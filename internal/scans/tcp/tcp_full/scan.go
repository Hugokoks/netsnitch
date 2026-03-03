package tcp_full

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/domain"
	"netsnitch/internal/fingerprint"
	"time"
)

func Scan(ctx context.Context, ip net.IP, port int, timeout time.Duration) domain.Result {
	result := domain.Result{
		Protocol: domain.TCP,
		IP:       ip,
		Port:     port,
		Open:     false,
	}

	// 1. Check if the port is open using a standard Dial
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return result // Connection failed, port is likely closed or filtered
	}

	// 2. Port is confirmed OPEN
	result.Open = true

	// 3. Delegate the complex service identification to the fingerprint engine.
	// We pass the initial 'conn' for the NullProbe phase.
	info := fingerprint.Identify(ctx, conn, ip, port, timeout)

	if info != nil {
		result.Service = info.Name
		result.Banner = info.Raw
	} else {
		result.Service = "unknown"
	}

	return result
}

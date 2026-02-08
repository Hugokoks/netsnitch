package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"netsnitch/internal/domain"
	"netsnitch/internal/probs"
)

func scanTarget(
	ctx context.Context,
	ip net.IP,
	port int,
	timeout time.Duration,
) domain.Result {

	address := fmt.Sprintf("%s:%d", ip, port)

	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return domain.Result{
			Protocol: domain.TCP,
			IP:       ip,
			Port:     port,
			Open:     false,
		}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	banner := ResolveBanner(conn)

	return domain.Result{
		Protocol: domain.TCP,
		IP:       ip,
		Port:     port,
		Open:     true,
		Banner:   banner,
		Service:  probs.DetectService(banner),
	}
}

package tcp_full

import (
	"context"
	"fmt"
	"net"
	"netsnitch/internal/domain"
	"time"
)

func Scan(ctx context.Context, ip net.IP, port int, timeout time.Duration) domain.Result {
	result := domain.Result{
		Protocol: domain.TCP,
		IP:       ip,
		Port:     port,
		Open:     false,
	}

	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return result
	}
	_ = conn.Close()

	result.Open = true
	return result
}

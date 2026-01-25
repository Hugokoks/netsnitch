package scanner

import (
	"context"
	"net"
	"netsnitch/internal/scan"
	"netsnitch/internal/target"
	"strconv"
	"time"
)

type TCPScanner struct {
	Timeout time.Duration
}


func (s *TCPScanner) Scan(ctx context.Context, t target.Target) scan.Result {
	result := scan.Result{
		IP:       t.IP,
		Port:     t.Port,
		Protocol: scan.TCP,
	}

	dialer := net.Dialer{Timeout: s.Timeout}
	addr := net.JoinHostPort(t.IP.String(), strconv.Itoa(t.Port))

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		result.Error = err
		return result
	}
	defer conn.Close()

	result.Open = true
	result.Banner = resolveBanner(conn, t)
	result.Service = detectService(result.Banner)

	return result
}
package tcp_full

import (
	"context"
	"fmt"
	"net"
	"time"

	"netsnitch/internal/domain"
	"netsnitch/internal/probs"
)

func Scan(
	ctx context.Context,
	ip net.IP,
	port int,
	timeout time.Duration,
) domain.Result {

	address := fmt.Sprintf("%s:%d", ip, port)

	///establish a connection to the port
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", address)

	result := domain.Result{
		Protocol: domain.TCP,
		IP:       ip,
		Port:     port,
		Open:     false,
	}

	if err != nil {
		return result
	}
	done := make(chan struct{})

	defer conn.Close()
	defer close(done)
	go endScan(ctx, done, conn)

	////set max read time
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	banner := ResolveBanner(conn)

	result.Open = true
	result.Banner = probs.DetectService(banner)
	return result
}

func endScan(ctx context.Context, done <-chan struct{}, conn net.Conn) {

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		case <-done:
			return
		}
	}
}

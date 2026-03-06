package fingerprint

import (
	"context"
	"fmt"
	"net"
	"time"
)

func (e *Engine) Identify(
	ctx context.Context,
	firstConn net.Conn,
	ip net.IP,
	port int,
	timeout time.Duration,
) *ServiceInfo {

	// Phase 1: passive banner
	raw := grabWithContext(ctx, firstConn, timeout)
	_ = firstConn.Close()

	if raw != "" {
		if info := e.Detect(port, raw); info != nil {
			return info
		}
	}

	// Phase 2: active probes
	probes := e.getProbes(port)
	address := fmt.Sprintf("%s:%d", ip, port)

	for _, p := range probes {

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		d := net.Dialer{Timeout: timeout}
		conn, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			continue
		}

		_ = conn.SetWriteDeadline(time.Now().Add(timeout))
		_ = conn.SetReadDeadline(time.Now().Add(timeout))

		// Send payload if any
		if len(p.Raw) > 0 {
			if _, err := conn.Write(p.Raw); err != nil {
				conn.Close()
				continue
			}
		}

		resp := grabWithContext(ctx, conn, timeout)
		conn.Close()

		if resp == "" {
			continue
		}

		if info := e.Detect(port, resp); info != nil {
			return info
		}
	}

	var info = &ServiceInfo{
		Service: "unknown",
		Banner:  raw,
	}
	return info
}

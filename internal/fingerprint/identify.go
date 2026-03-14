package fingerprint

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Identify is used when you DO NOT already have an open TCP connection.
// Typical use: stealth SYN scan found port as OPEN, now do normal connect for fingerprinting.
func (e *Engine) Identify(
	ctx context.Context,
	ip net.IP,
	port int,
	timeout time.Duration,
) *ServiceInfo {
	address := fmt.Sprintf("%s:%d", ip.String(), port)

	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil
	}

	return e.IdentifyWithConn(ctx, conn, ip, port, timeout)
}

// IdentifyWithConn is used when you ALREADY have an open TCP connection.
// Typical use: full TCP connect scan.
func (e *Engine) IdentifyWithConn(
	ctx context.Context,
	firstConn net.Conn,
	ip net.IP,
	port int,
	timeout time.Duration,
) *ServiceInfo {
	if firstConn == nil {
		return nil
	}

	// Keep last useful response for unknown fallback
	lastResp := ""

	// Phase 1: passive banner
	raw := grabWithContext(ctx, firstConn, timeout)
	_ = firstConn.Close()

	if raw != "" {
		lastResp = raw

		if info := e.Detect(port, raw); info != nil {
			return info
		}
	}

	// Phase 2: active probes
	probes := e.getProbes(port)
	address := fmt.Sprintf("%s:%d", ip.String(), port)

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

		if len(p.Raw) > 0 {
			if _, err := conn.Write(p.Raw); err != nil {
				_ = conn.Close()
				continue
			}
		}

		resp := grabWithContext(ctx, conn, timeout)
		_ = conn.Close()

		if resp == "" {
			continue
		}

		lastResp = resp

		if info := e.Detect(port, resp); info != nil {
			return info
		}
	}

	return &ServiceInfo{
		Service: "unknown",
		Banner:  lastResp,
	}
}

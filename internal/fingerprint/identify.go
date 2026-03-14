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
		fmt.Printf("[fp-passive] port=%d raw=%q\n", port, raw)
		if info := e.Detect(port, raw); info != nil {
			fmt.Printf("[fp-hit-passive] port=%d service=%s rule=%s\n", port, info.Service, info.RuleID)
			return info
		}
		fmt.Printf("[fp-miss-passive] port=%d raw=%q\n", port, raw)
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
		fmt.Printf("[fp-probe-start] port=%d probe=%s payload_len=%d\n", port, p.ID, len(p.Raw))
		d := net.Dialer{Timeout: timeout}
		conn, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			fmt.Printf("[fp-probe-dial-fail] port=%d probe=%s err=%v\n", port, p.ID, err)
			continue
		}

		_ = conn.SetWriteDeadline(time.Now().Add(timeout))
		_ = conn.SetReadDeadline(time.Now().Add(timeout))

		if len(p.Raw) > 0 {
			if _, err := conn.Write(p.Raw); err != nil {
				fmt.Printf("[fp-probe-write-fail] port=%d probe=%s err=%v\n", port, p.ID, err)
				_ = conn.Close()
				continue
			}
			fmt.Printf("[fp-probe-write-ok] port=%d probe=%s\n", port, p.ID)
		}

		resp := grabWithContext(ctx, conn, timeout)
		_ = conn.Close()
		fmt.Printf("[fp-probe-resp] port=%d probe=%s resp_len=%d resp=%q\n", port, p.ID, len(resp), resp)
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

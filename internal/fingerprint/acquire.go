package fingerprint

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Identify tries to determine the service by using passive and active probing
func Identify(ctx context.Context, firstConn net.Conn, ip net.IP, port int, timeout time.Duration) *ServiceInfo {
	// --- PHASE 1: Passive Null Probe ---
	// Use the already established connection to wait for an initial banner (e.g., SSH, FTP)
	raw := grabWithContext(ctx, firstConn, timeout)
	firstConn.Close() // Always close the initial connection after reading

	if raw != "" {
		return Detect(raw) // Match against Recog XML patterns
	}

	// --- PHASE 2: Active Probing ---
	// If the server is silent, we start sending specific payloads
	probes := GetProbesForPort(port)
	address := fmt.Sprintf("%s:%d", ip, port)

	for _, p := range probes {
		// Check if the scan was cancelled (e.g., Ctrl+C)
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Skip NullProbe as we already handled it in Phase 1
		if p.Name == "NullProbe" {
			continue
		}

		// Dial a fresh connection for each probe to ensure a clean protocol state
		d := net.Dialer{Timeout: timeout}
		conn, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			continue // Port might have become unreachable or rate-limited
		}

		// Send the active probe payload
		conn.SetWriteDeadline(time.Now().Add(timeout))
		_, err = conn.Write(p.RawData)
		if err != nil {
			conn.Close()
			continue
		}

		// Wait for a response to our probe
		response := grabWithContext(ctx, conn, timeout)
		conn.Close()

		if response != "" {
			info := Detect(response)
			// If we successfully identified the product, return immediately
			if info != nil && info.Name != "unknown-text" {
				return info
			}
		}
	}

	return nil
}

// grabWithContext reads data from a connection while respecting the provided context
func grabWithContext(ctx context.Context, conn net.Conn, timeout time.Duration) string {
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(timeout))

	type readResult struct {
		n   int
		err error
	}
	resChan := make(chan readResult, 1)

	// Perform the read in a goroutine so we can select on the context
	go func() {
		n, err := conn.Read(buf)
		resChan <- readResult{n, err}
	}()

	select {
	case <-ctx.Done():
		return ""
	case res := <-resChan:
		if res.err != nil || res.n == 0 {
			return ""
		}
		return string(buf[:res.n])
	}
}

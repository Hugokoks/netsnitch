package scanner

import (
	"net"
	"netsnitch/internal/scanner/probs"
	"netsnitch/internal/target"
	"strings"
	"time"
)

func grabBanner(conn net.Conn) string {
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	
	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
}

func resolveBanner(conn net.Conn, t target.Target) string {
	// 1.banner
	if banner := grabBanner(conn); banner != "" {
		return banner
	}

	// 2.HTTP probe (port-agnostic)
	if banner := probs.TryHTTP(conn, t.IP.String()); banner != "" {
		return banner
	}

	return ""
}
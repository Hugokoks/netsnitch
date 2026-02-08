package tcp

import (
	"net"
	"netsnitch/internal/probs"
	"strings"
)

func GrabBanner(conn net.Conn) string {

	buf := make([]byte, 512)
	n, err := conn.Read(buf)

	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
}

func ResolveBanner(conn net.Conn) string {
	// 1.banner
	if banner := GrabBanner(conn); banner != "" {
		return banner
	}

	// 2.HTTP probe (port-agnostic)
	if banner := probs.TryHTTP(conn); banner != "" {
		return banner
	}

	return ""
}

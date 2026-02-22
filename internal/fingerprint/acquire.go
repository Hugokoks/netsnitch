package fingerprint

import (
	"net"
	"time"
)

func AcquireTCP(conn net.Conn, port int, timeout time.Duration) string {

	raw := grabBanner(conn, timeout)

	if raw == "" && looksLikeHTTP(port) {
		raw = probeHTTP(conn, timeout)
	}

	return raw
}

func grabBanner(conn net.Conn, timeout time.Duration) string {

	buf := make([]byte, 1024)

	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ""
	}

	return string(buf[:n])
}

func looksLikeHTTP(port int) bool {
	switch port {
	case 80, 8080, 8180, 8000:
		return true
	default:
		return false
	}
}

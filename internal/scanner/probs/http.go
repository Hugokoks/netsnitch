package probs

import (
	"net"
	"strings"
	"time"
)

func TryHTTP(conn net.Conn, host string) string {
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))

	req := "HEAD / HTTP/1.1\r\nHost: " + host + "\r\n\r\n"
	_, _ = conn.Write([]byte(req))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
}

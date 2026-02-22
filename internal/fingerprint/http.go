package fingerprint

import (
	"net"
	"strings"
	"time"
)

func probeHTTP(conn net.Conn, timeout time.Duration) string {

	req := "HEAD / HTTP/1.0\r\n\r\n"

	conn.Write([]byte(req))

	buf := make([]byte, 2048)
	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	n, _ := conn.Read(buf)
	if n == 0 {
		return ""
	}

	return string(buf[:n])
}

func parseHTTP(raw string) *ServiceInfo {

	info := &ServiceInfo{
		Name: "http",
		Raw:  raw,
	}

	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "server:") {
			info.Product = strings.TrimSpace(line)
			break
		}
	}

	return info
}

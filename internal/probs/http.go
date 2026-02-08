package probs

import (
	"net"
	"strings"
)

func TryHTTP(conn net.Conn) string {

	////check data type of addr
	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return ""
	}

	host := addr.IP.String()

	req := "HEAD / HTTP/1.1\r\n" +
		"Host: " + host + "\r\n" +
		"Connection: close\r\n\r\n"

	_, _ = conn.Write([]byte(req))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
}

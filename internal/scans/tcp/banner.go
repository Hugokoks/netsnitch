package tcp

import (
	"net"
	"strings"
	"time"
)

func GrabBanner(conn net.Conn) string {

	buf := make([]byte, 512)

	n, err := conn.Read(buf)

	if err != nil || n == 0 {
		return ""
	}

	return strings.TrimSpace(string(buf[:n]))
}
func tryHTTP(conn net.Conn) string {

	request := "GET / HTTP/1.1\r\nHost: " + conn.RemoteAddr().String() + "\r\nConnection: close\r\n\r\n"

	conn.Write([]byte(request))

	htmlLenght := 512

	buf := make([]byte, htmlLenght)

	_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	n, _ := conn.Read(buf)
	return string(buf[:n])
}

func ResolveBanner(conn net.Conn) string {
	// 1.banner
	if banner := GrabBanner(conn); banner != "" {
		return banner
	}

	// 2.HTTP probe (port-agnostic)
	if banner := tryHTTP(conn); banner != "" {
		return banner
	}

	return ""
}

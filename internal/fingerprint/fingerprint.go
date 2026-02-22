package fingerprint

import (
	"net"
	"time"
)

type ServiceInfo struct {
	Name    string
	Product string
	Version string
	Raw     string
}

func Detect(conn net.Conn, port int, timeout time.Duration) *ServiceInfo {

	raw := acquireRaw(conn, port, timeout)

	if raw == "" {
		return &ServiceInfo{Name: "unknown"}
	}

	return matchPatterns(raw)
}

func matchPatterns(raw string) *ServiceInfo {

	for _, pattern := range registry {

		if matches := pattern.Regex.FindStringSubmatch(raw); matches != nil {
			info := pattern.Parser(matches)
			info.Raw = raw
			return info
		}
	}

	return &ServiceInfo{
		Name: "unknown-text",
		Raw:  raw,
	}
}

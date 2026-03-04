package fingerprint

import "strings"

func guessProtocol(banner string) string {

	s := strings.ToLower(banner)

	switch {

	case strings.Contains(s, "ssh"):
		return "ssh"

	case strings.Contains(s, "ftp"):
		return "ftp"

	case strings.Contains(s, "smtp"):
		return "smtp"

	case strings.Contains(s, "mysql"):
		return "mysql"

	case strings.Contains(s, "http"):
		return "http"

	case strings.Contains(s, "redis"):
		return "redis"

	}

	return ""
}

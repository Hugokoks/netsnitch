package probs

import (
	"strings"
)

func DetectService(banner string) string {
	b := strings.ToLower(banner)

	switch {
	case strings.HasPrefix(b, "ssh-"):
		return "ssh"

	case strings.Contains(b, "http/"):
		return "http"

	case strings.HasPrefix(b, "220"):
		return "smtp/ftp"

	case strings.Contains(b, "mysql"):
		return "mysql"

	case b != "":
		return "unknown-text"

	default:
		return "unknown"
	}
}




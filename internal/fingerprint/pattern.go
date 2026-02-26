package fingerprint

import "regexp"

type Pattern struct {
	Service string
	Regex   *regexp.Regexp
	Parser  func([]string) *ServiceInfo
}

var registry = []Pattern{
	{
		Service: "ssh",
		Regex:   regexp.MustCompile(`^SSH-\d\.\d-([^\r\n]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "ssh",
				Product: matches[1],
			}
		},
	},
	{
		Service: "ftp",
		Regex:   regexp.MustCompile(`vsftpd\s+([\d\.]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "ftp",
				Product: "vsftpd",
				Version: matches[1],
			}
		},
	},
	{
		Service: "ftp",
		Regex:   regexp.MustCompile(`ProFTPD\s+([\d\.]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "ftp",
				Product: "ProFTPD",
				Version: matches[1],
			}
		},
	},
	{
		Service: "smtp",
		Regex:   regexp.MustCompile(`ESMTP\s+([^\s]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "smtp",
				Product: matches[1],
			}
		},
	},
	{
		Service: "http",
		Regex:   regexp.MustCompile(`Server:\s*([^\r\n]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "http",
				Product: matches[1],
			}
		},
	},
	{
		Service: "mysql",
		Regex:   regexp.MustCompile(`(?i)[\x00-\xff]{3,5}([58]\.[0-9]+\.[0-9]+[^\x00-\x1f]+)`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "mysql",
				Product: "MySQL",
				Version: matches[1],
			}
		},
	},
	{
		Service: "redis",
		Regex:   regexp.MustCompile(`(?i)-ERR unknown command|PONG`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "redis",
				Product: "Redis",
			}
		},
	},
	{
		Service: "postgresql",
		Regex:   regexp.MustCompile(`(?i)PostgreSQL`),
		Parser: func(matches []string) *ServiceInfo {
			return &ServiceInfo{
				Name:    "postgresql",
				Product: "PostgreSQL",
			}
		},
	},
}

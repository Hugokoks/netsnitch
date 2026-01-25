package target

import (
	"net"
	"netsnitch/internal/scan"
)

type Builder interface {
	Protocol() scan.Protocol
	Build(ips []net.IP) []Target
}

var builders = map[scan.Protocol]Builder{}

func Register(b Builder) {
	builders[b.Protocol()] = b
}

func BuildTargets(cfg scan.Config, ips []net.IP) []Target {
	if b, ok := builders[cfg.Type]; ok {
		return b.Build(ips)
	}

	panic("no target builder registered for scan type")
}

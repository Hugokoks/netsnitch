package tasks

import (
	"netsnitch/internal/domain"
)

type Builder interface {
	Protocol() domain.Protocol
	Build(cidr string, cfg domain.Config) []Task
}

var builders = map[domain.Protocol]Builder{}

func Register(b Builder) {
	builders[b.Protocol()] = b
}

func Build(cfg domain.Config, cidr string) []Task {
	if b, ok := builders[cfg.Type]; ok {
		return b.Build(cidr, cfg)
	}

	panic("no task builder registered for scan type")
}
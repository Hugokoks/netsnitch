package tasks

import (
	"netsnitch/internal/domain"
)

type Builder interface {
	Protocol() domain.Protocol
	Build(cfg domain.Config) []Task
}

var builders = map[domain.Protocol]Builder{}

func Register(b Builder) {
	builders[b.Protocol()] = b
}

func Build(cfg domain.Config) []Task {
	if b, ok := builders[cfg.Type]; ok {
		return b.Build(cfg)
	}

	panic("no task builder registered for scan type")
}

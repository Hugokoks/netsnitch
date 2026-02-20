package udp

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.UDP
}
func (b Builder) Build(cfg domain.Config) []tasks.Task {
	return []tasks.Task{}
}

func init() {

	tasks.Register(Builder{})
}

package arp_active

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.ARP_ACTIVE
}

func (b Builder) Build(cfg domain.Config) []tasks.Task {
	return []tasks.Task{
		&Task{
			timeout: cfg.Timeout,
			scope:   cfg.Scope,
		},
	}
}

func init() {
	tasks.Register(Builder{})
}

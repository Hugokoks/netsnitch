package arp_active

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.ARP
}

func (b Builder) Build(cfg domain.Config) []tasks.Task {
	return []tasks.Task{
		&Task{
			timeout: cfg.Timeout,
			scope:   cfg.Scope,
			mode:    cfg.Mode,
			render:  cfg.Render,
		},
	}
}

func init() {
	tasks.Register(Builder{})
}

package arp_active

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {

	return domain.ARP_ACTIVE
}

func (b Builder) Build(cidr string, cfg domain.Config) []tasks.Task {

	var tasks []tasks.Task

	tasks = append(tasks, &Task{timeout: cfg.Timeout, cidr: cidr})

	return tasks

}

func init() {
	tasks.Register(Builder{})
}

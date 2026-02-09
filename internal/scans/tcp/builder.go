package tcp

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.TCP
}

func (b Builder) Build(cfg domain.Config) []tasks.Task {

	ips, err := domain.ResolveScope(cfg.Scope)
	if err != nil {
		panic(err)
	}

	ports, err := domain.ResolvePortScope(cfg.Ports)
	if err != nil {

		panic(err)

	}

	var tasks []tasks.Task

	for _, ip := range ips {
		for _, port := range ports {
			tasks = append(tasks, &Task{
				ip:      ip,
				port:    port,
				timeout: cfg.Timeout,
				mode:    cfg.Mode,
			})
		}
	}

	return tasks
}

func init() {
	tasks.Register(Builder{})
}

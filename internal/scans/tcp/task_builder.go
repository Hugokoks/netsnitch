package tcp

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.TCP
}

func (b Builder) Build(cidr string, cfg domain.Config) []tasks.Task {
	ips, err := domain.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	var tasks []tasks.Task

	for _, ip := range ips {
		for _, port := range domain.DefaultPorts {
			tasks = append(
				tasks,
				&Task{ip: ip, port: port, timeout: cfg.Timeout},
			)
		}
	}

	return tasks
}

func init() {
	tasks.Register(Builder{})
}

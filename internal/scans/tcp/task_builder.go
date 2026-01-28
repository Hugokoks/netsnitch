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

	var tasksStack []tasks.Task

	for _, ip := range ips {
		for _, port := range domain.DefaultPorts {
			tasksStack = append(
				tasksStack,
				NewTask(ip, port, cfg.Timeout),
			)
		}
	}

	return tasksStack
}

func init() {
	tasks.Register(Builder{})
}

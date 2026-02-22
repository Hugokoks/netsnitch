package udp

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/udp/udp_scan"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.UDP
}
func (b Builder) Build(cfg domain.Config) ([]tasks.Task, error) {

	ips, err := domain.ResolveScope(cfg.Scope)
	if err != nil {
		return nil, fmt.Errorf("UDP builder Resolve scope error %w", err)
	}

	ports, err := domain.ResolvePortScope(cfg.Ports)

	if err != nil {

		return nil, fmt.Errorf("UDP builder Resolve ports error %w", err)

	}

	// Create UDP socket manager once

	mgr, err := udp_scan.NewManager()
	if err != nil {
		return nil, fmt.Errorf("UDP builder create manager error %w", err)
	}

	var taskList []tasks.Task
	for _, ip := range ips {
		for _, port := range ports {
			task := UDPTask{
				ip: ip, port: port, timeout: cfg.Timeout,
				render: cfg.Render, openOnly: cfg.OpenOnly,
				mgr: mgr,
			}
			taskList = append(taskList, &task)
		}
	}

	return taskList, nil
}

func init() {

	tasks.Register(Builder{})
}

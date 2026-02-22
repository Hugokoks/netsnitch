package tcp

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/tcp/tcp_stealth"
	"netsnitch/internal/tasks"
)

type Builder struct{}

func (b Builder) Protocol() domain.Protocol {
	return domain.TCP
}

func (b Builder) Build(cfg domain.Config) ([]tasks.Task, error) {

	ips, err := domain.ResolveScope(cfg.Scope)
	if err != nil {
		return nil, fmt.Errorf("TCP builder Resolve scope error %w", err)
	}

	ports, err := domain.ResolvePortScope(cfg.Ports)
	if err != nil {
		return nil, fmt.Errorf("TCP builder Resolve scope error %w", err)
	}

	////Open network socket only onc
	var mgr *tcp_stealth.Manager
	if cfg.Mode == domain.STEALTH {
		mgr, err = tcp_stealth.NewManager()
		if err != nil {
			return nil, fmt.Errorf("TCP builder create manager error %w", err)
		}
	}

	var taskList []tasks.Task
	for _, ip := range ips {
		for _, port := range ports {
			base := baseTask{
				ip: ip, port: port, timeout: cfg.Timeout,
				render: cfg.Render, openOnly: cfg.OpenOnly,
			}
			if cfg.Mode == domain.STEALTH {
				taskList = append(taskList, &StealthTask{baseTask: base, mgr: mgr})
			} else {
				taskList = append(taskList, &FullTask{baseTask: base})
			}
		}
	}
	return taskList, nil
}

func init() {
	tasks.Register(Builder{})
}

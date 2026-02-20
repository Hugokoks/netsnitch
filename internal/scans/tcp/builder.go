package tcp

import (
	"netsnitch/internal/domain"
	"netsnitch/internal/scans/tcp/tcp_stealth"
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

	////Open network socket only onc
	var mgr *tcp_stealth.Manager
	if cfg.Mode == domain.STEALTH {
		mgr, err = tcp_stealth.NewManager()
		if err != nil {
			panic(err)
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
	return taskList
}

func init() {
	tasks.Register(Builder{})
}

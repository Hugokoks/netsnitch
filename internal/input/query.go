package input

import (
	"netsnitch/internal/domain"
)

type Query struct {
	Stages []Stage
}

type Stage struct {
	Protocol domain.Protocol ////ARP_ACTIVE,TCP,
	Ports    domain.PortScope
	Scope    domain.Scope
	////IP scope single ip 192,168.0.1,cidr 192.168.0.0/24, multiple ips 192.168.0.1,192.168.0.02
	Concurrency int ////how many workers will do scan
}

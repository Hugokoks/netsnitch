package domain

import (
	"time"
)

type Config struct {
	Type    Protocol
	Timeout time.Duration

	Scope Scope ////scope of ip's cidr, ip's,single ip
	Ports PortScope
}

var DefaultPorts = []int{22, 80, 443}

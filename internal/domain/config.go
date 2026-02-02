package domain

import (
	"time"
)

type Config struct {
	Type        Protocol
	Timeout     time.Duration
	Concurrency int

	Scope Scope ////scope of ip's cidr, ip's,single ip
	Ports []int
}

var DefaultPorts = []int{22, 80, 443}

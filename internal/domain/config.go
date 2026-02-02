package domain

import (
	"time"
)

type Config struct {
	Type        Protocol
	Timeout     time.Duration
	Concurrency int

	Scope   Scope
	Options []int
}

var DefaultPorts = []int{22, 80, 443}

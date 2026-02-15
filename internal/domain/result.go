package domain

import (
	"net"
	"time"
)

type Result struct {
	IP       net.IP
	MAC      net.HardwareAddr
	Port     int
	Protocol Protocol
	OutputType

	Open  bool
	Alive bool
	RTT   time.Duration

	Banner  string
	Service string
	Error   error
}

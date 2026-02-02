package domain

import (
	"net"
	"time"
)

type Result struct {
	IP       net.IP
	MAC net.HardwareAddr
	Port     int
	Protocol Protocol

	Open     bool    
	Alive   bool     
	RTT      time.Duration

	Banner   string
	Service  string
	Error    error
}
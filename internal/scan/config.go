package scan

import "time"

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
	ICMP Protocol = "icmp"
	ARP Protocol = "arp"
)

type Config struct {
	Type        Protocol
	Timeout     time.Duration
	Concurrency int
}

var DefaultPorts = []int{22, 80, 443}

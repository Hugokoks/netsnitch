package domain

import "time"

type Protocol string

const (
	TCP         Protocol = "tcp"
	UDP         Protocol = "udp"
	ICMP        Protocol = "icmp"
	ARP_PASSIVE Protocol = "arp_passive"
	ARP_ACTIVE  Protocol = "arp_active"
)

type Config struct {
	Type        Protocol
	Timeout     time.Duration
	Concurrency int
}

var DefaultPorts = []int{22, 80, 443}

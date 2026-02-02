package domain

import (
	"fmt"
	"strings"
)

type Protocol string

const (
	TCP         Protocol = "tcp"
	UDP         Protocol = "udp"
	ICMP        Protocol = "icmp"
	ARP_PASSIVE Protocol = "arp_passive"
	ARP_ACTIVE  Protocol = "arp_active"
)

func ParseProtocol(s string) (Protocol, error) {
	switch strings.ToLower(s) {
	case "tcp":
		return TCP, nil
	case "arp":
		return ARP_ACTIVE, nil
	default:
		return "", fmt.Errorf("unknown protocol: %s", s)
	}
}

package domain

type Protocol string

const (
	TCP         Protocol = "tcp"
	UDP         Protocol = "udp"
	ICMP        Protocol = "icmp"
	ARP_PASSIVE Protocol = "arp_passive"
	ARP_ACTIVE  Protocol = "arp_active"
)

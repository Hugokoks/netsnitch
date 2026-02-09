package domain

type Protocol string
type ScanMode string

const (
	TCP  Protocol = "tcp"
	UDP  Protocol = "udp"
	ICMP Protocol = "icmp"
	ARP  Protocol = "arp"

	FULL        ScanMode = "f"
	STEALH      ScanMode = "s"
	ARP_ACTIVE  ScanMode = "active"
	ARP_PASSIVE ScanMode = "passive"
)

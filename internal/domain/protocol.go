package domain

type Protocol string
type ScanMode string
type OutputType string

const (
	TCP  Protocol = "tcp"
	UDP  Protocol = "udp"
	ICMP Protocol = "icmp"
	ARP  Protocol = "arp"

	FULL        ScanMode = "f"
	STEALTH     ScanMode = "s"
	ARP_ACTIVE  ScanMode = "active"
	ARP_PASSIVE ScanMode = "passive"

	ROWS_OUT OutputType = "rows"
	JSON_OUT OutputType = "json"
)

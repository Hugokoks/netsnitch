package domain

type Protocol string
type ScanMode string
type RenderType string

const (
	TCP  Protocol = "tcp"
	UDP  Protocol = "udp"
	ICMP Protocol = "icmp"
	ARP  Protocol = "arp"

	FULL        ScanMode = "f"
	STEALTH     ScanMode = "s"
	ARP_ACTIVE  ScanMode = "active"
	ARP_PASSIVE ScanMode = "passive"

	ROWS_OUT RenderType = "rows"
	JSON_OUT RenderType = "json"
)

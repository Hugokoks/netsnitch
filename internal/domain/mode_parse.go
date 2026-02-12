package domain

import "fmt"

func ParseScanMode(s string) (ScanMode, error) {

	switch s {

	case string(FULL):
		return FULL, nil

	case string(STEALTH):
		return STEALTH, nil

	case string(ARP_ACTIVE):
		return ARP_ACTIVE, nil

	case string(ARP_PASSIVE):
		return ARP_PASSIVE, nil

	default:
		return "", fmt.Errorf("invalid scan mode: %s", s)
	}
}

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

func ParseRenderType(s string) (RenderType, error) {

	switch s {

	case string(JSON_OUT):
		return JSON_OUT, nil

	case string(ROWS_OUT):
		return ROWS_OUT, nil
	default:
		return "", fmt.Errorf("invalid render type %s", s)

	}

}

package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/input"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {

	return domain.ARP

}

func (Parser) Parse(cfg *domain.Config, rest []string, flags input.Flags) error {

	if len(rest) < 2 {
		return fmt.Errorf("usage: arp <cidr>")
	}

	//// ---- scope ----
	scope, err := domain.ParseScope(rest[1])
	if err != nil {
		return err
	}

	// ----apply settings ----
	cfg.Scope = scope

	return nil
}
func (Parser) ApplyDefaults(cfg *domain.Config) {

	if cfg.Timeout <= 0 {
		cfg.Timeout = domain.DefaultTimeout
	}

	if cfg.Mode == "" {
		cfg.Mode = domain.ARP_ACTIVE
	}

}

func init() {

	input.Register(Parser{})
}

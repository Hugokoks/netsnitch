package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/input"
	"time"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {

	return domain.ARP

}

func (Parser) Parse(tokens []string) (domain.Config, error) {

	if len(tokens) < 2 {
		return domain.Config{}, fmt.Errorf("usage: arp <cidr>")
	}

	flags, rest, err := input.ExtractFlags(tokens[1:])

	if err != nil {
		return domain.Config{}, err
	}

	if len(rest) != 1 {
		return domain.Config{}, fmt.Errorf("exactly one target scope required")
	}

	//// ---- scope ----
	scope, err := domain.ParseScope(rest[0])
	if err != nil {
		return domain.Config{}, err
	}

	// ----timeout-----
	var timeout time.Duration

	if t, ok := flags["t"]; ok {
		if d, err := time.ParseDuration(t); err == nil {
			timeout = d
		}
	}

	return domain.Config{
		Type:    domain.ARP,
		Scope:   scope,
		Timeout: timeout,
	}, nil
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

package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/input"
	"time"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {

	return domain.ARP_ACTIVE

}

func (Parser) Parse(tokens []string) (input.Stage, error) {

	if len(tokens) < 2 {
		return input.Stage{}, fmt.Errorf("usage: arp <cidr>")
	}

	flags, rest, err := input.ExtractFlags(tokens[1:])

	if err != nil {
		return input.Stage{}, err
	}

	if len(rest) != 1 {
		return input.Stage{}, fmt.Errorf("exactly one target scope required")
	}

	//// ---- scope ----
	scope, err := domain.ParseScope(rest[0])
	if err != nil {
		return input.Stage{}, err
	}

	// ----timeout-----
	timeout := domain.DefaultTimeout

	if t, ok := flags["t"]; ok {
		if d, err := time.ParseDuration(t); err == nil {
			timeout = d
		}
	}

	return input.Stage{
		Protocol: domain.ARP_ACTIVE,
		Scope:    scope,
		Timeout:  timeout,
	}, nil
}

func init() {

	input.Register(Parser{})
}

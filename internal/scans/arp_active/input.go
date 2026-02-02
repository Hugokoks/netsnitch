package arp_active

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/input"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {

	return domain.ARP_ACTIVE

}

func (Parser) Parse(tokens []string) (input.Stage, error) {

	if len(tokens) != 2 {

		return input.Stage{}, fmt.Errorf("usage: arp <cidr>")

	}

	return input.Stage{
		Protocol: domain.ARP_ACTIVE,
		Scope: domain.Scope{
			Type: domain.ScopeCIDR,
			CIDR: tokens[1],
		},
	}, nil

}

func init() {

	input.Register(Parser{})
}

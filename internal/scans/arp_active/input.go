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
		return input.Stage{}, fmt.Errorf("usage: arp <cidr|ip[,ip]>")
	}

	scope, err := input.ParseScope(tokens[1])

	if err != nil {
		return input.Stage{}, err
	}

	return input.Stage{
		Protocol: domain.ARP_ACTIVE,
		Scope:    scope,
	}, nil
}

func init() {

	input.Register(Parser{})
}

package tcp

import (
	"fmt"

	"netsnitch/internal/domain"
	"netsnitch/internal/input"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {
	return domain.TCP
}

func (Parser) Parse(tokens []string) (input.Stage, error) {
	if len(tokens) < 2 {
		return input.Stage{}, fmt.Errorf(
			"usage: tcp [--ports:<ports>] <cidr|ip>",
		)
	}

	flags, rest, err := input.ExtractFlags(tokens[1:])
	if err != nil {
		return input.Stage{}, err
	}

	if len(rest) != 1 {
		return input.Stage{}, fmt.Errorf("exactly one target scope required")
	}

	// ----scope ----
	scope, err := domain.ParseScope(rest[0])
	if err != nil {
		return input.Stage{}, err
	}

	// ----ports----
	portScope := domain.PortScope{
		Type:  domain.PortsList,
		Ports: domain.DefaultPorts,
	}

	if p, ok := flags["ports"]; ok {
		ps, err := domain.ParsePortScope(p)
		if err != nil {
			return input.Stage{}, err
		}
		portScope = ps
	}

	return input.Stage{
		Protocol: domain.TCP,
		Scope:    scope,
		Ports:    portScope,
	}, nil
}

func init() {
	input.Register(Parser{})
}

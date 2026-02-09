package tcp

import (
	"fmt"
	"time"

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
			"usage: tcp [--p <p>] <cidr|ip>",
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

	if p, ok := flags["p"]; ok {
		ps, err := domain.ParsePortScope(p)

		if err != nil {
			return input.Stage{}, err
		}

		portScope = ps
	}

	// ----timeout-----
	timeout := domain.DefaultTimeout

	if t, ok := flags["t"]; ok {
		if d, err := time.ParseDuration(t); err == nil {
			timeout = d
		}
	}

	return input.Stage{
		Protocol: domain.TCP,
		Scope:    scope,
		Ports:    portScope,
		Timeout:  timeout,
	}, nil
}

func init() {
	input.Register(Parser{})
}

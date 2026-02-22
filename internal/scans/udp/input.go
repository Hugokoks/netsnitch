package udp

import (
	"fmt"
	"netsnitch/internal/domain"
	"netsnitch/internal/input"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {
	return domain.UDP
}

func (Parser) Parse(cfg *domain.Config, rest []string, flags input.Flags) error {

	if len(rest) < 2 {
		return fmt.Errorf("usage: udp [--p <p>] <cidr|ip>")
	}
	ipToken := rest[len(rest)-1]

	// ----scope----
	scope, err := domain.ParseScope(ipToken)
	if err != nil {
		return err
	}

	// ----ports----
	var portScope domain.PortScope
	if p, ok := flags["p"]; ok {
		portScope, _ = domain.ParsePortScope(p)

	}

	// ----OpenOnly----
	if _, ok := flags["open"]; ok {

		cfg.OpenOnly = true

	}

	cfg.Ports = portScope
	cfg.Scope = scope

	return nil

}

func (Parser) ApplyDefaults(cfg *domain.Config) {

	if len(cfg.Ports.Ports) == 0 && cfg.Ports.Type != domain.PortsRange {
		cfg.Ports = domain.PortScope{
			Type:  domain.PortsList,
			Ports: domain.DefaultPorts,
		}
	}

}

func init() {

	input.Register(Parser{})
}

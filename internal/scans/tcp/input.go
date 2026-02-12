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

func (Parser) Parse(tokens []string) (domain.Config, error) {
	if len(tokens) < 2 {
		return domain.Config{}, fmt.Errorf(
			"usage: tcp [--p <p>] <cidr|ip>",
		)
	}

	flags, rest, err := input.ExtractFlags(tokens[1:])
	if err != nil {
		return domain.Config{}, err
	}

	if len(rest) != 1 {
		return domain.Config{}, fmt.Errorf("exactly one target scope required")
	}

	// ----scope ----
	scope, err := domain.ParseScope(rest[0])
	if err != nil {
		return domain.Config{}, err
	}

	// ----ports----
	var portScope domain.PortScope
	if p, ok := flags["p"]; ok {
		ps, err := domain.ParsePortScope(p)
		if err != nil {
			return domain.Config{}, err
		}

		portScope = ps
	}

	// ----timeout-----
	var timeout time.Duration

	if t, ok := flags["t"]; ok {
		if d, err := time.ParseDuration(t); err == nil {
			timeout = d
		}
	}

	// ----mode----

	var mode domain.ScanMode
	if m, ok := flags["mode"]; ok {

		parseMode, err := domain.ParseScanMode(m)
		if err != nil {

			return domain.Config{}, fmt.Errorf("mode %s doesn't exist", m)

		}
		mode = parseMode

	}

	return domain.Config{
		Type:    domain.TCP,
		Scope:   scope,
		Ports:   portScope,
		Timeout: timeout,
		Mode:    mode,
	}, nil

}

func (Parser) ApplyDefaults(cfg *domain.Config) {

	if cfg.Timeout <= 0 {
		cfg.Timeout = domain.DefaultTimeout
	}
	if cfg.Mode == "" {
		cfg.Mode = domain.FULL
	}
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

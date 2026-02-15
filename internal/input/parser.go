package input

import (
	"fmt"
	"netsnitch/internal/domain"
)

type Parser interface {
	Protocol() domain.Protocol
	Parse(cfg *domain.Config, rest []string, flags Flags) error
	ApplyDefaults(cfg *domain.Config)
}

var parsers = map[domain.Protocol]Parser{}

func Register(p Parser) {
	parsers[p.Protocol()] = p
}

func getParser(proto domain.Protocol) (Parser, error) {
	p, ok := parsers[proto]
	if !ok {
		return nil, fmt.Errorf("no input parser for protocol %s", proto)
	}
	return p, nil
}

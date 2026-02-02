package tcp

import (
	"fmt"
	"strconv"
	"strings"

	"netsnitch/internal/domain"
	"netsnitch/internal/input"
)

type Parser struct{}

func (Parser) Protocol() domain.Protocol {
	return domain.TCP
}
func (Parser) Parse(tokens []string) (input.Stage, error) {
	var (
		ports []int
		scope domain.Scope
		err   error
	)

	switch len(tokens) {
	case 2:
		// tcp <scope>
		ports = domain.DefaultPorts
		scope, err = input.ParseScope(tokens[1])

	case 3:
		// tcp <ports> <scope>
		ports, err = parsePorts(tokens[1])
		if err != nil {
			return input.Stage{}, err
		}
		scope, err = input.ParseScope(tokens[2])

	default:
		return input.Stage{}, fmt.Errorf("usage: tcp [ports] <cidr|ip>")
	}

	if err != nil {
		return input.Stage{}, err
	}

	return input.Stage{
		Protocol: domain.TCP,
		Scope:    scope,
		Ports:    ports,
	}, nil
}

func parsePorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ports := make([]int, 0, len(parts))

	for _, p := range parts {
		port, err := strconv.Atoi(p)
		if err != nil || port <= 0 || port > 65535 {
			return nil, fmt.Errorf("invalid port: %s", p)
		}
		ports = append(ports, port)
	}

	return ports, nil
}

func init() {
	input.Register(Parser{})
}

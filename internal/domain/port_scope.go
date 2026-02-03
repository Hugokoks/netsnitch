package domain

import (
	"fmt"
	"strconv"
	"strings"
)

/////Port flag
/////--port:5 ///single port
/////--port:5-10 ////port range
/////--port:1,3,44,66,5000 ////list of ports
/////--port:all ////all ports

type PortScopeType int

const (
	PortSingle PortScopeType = iota
	PortsList
	PortsRange
	PortsAll
)

type PortScope struct {
	Type  PortScopeType
	Ports []int
	From  int
	To    int
}

func ResolvePortScope(scope PortScope) ([]int, error) {
	switch scope.Type {
	case PortSingle, PortsList:
		return scope.Ports, nil

	case PortsRange:
		var ports []int
		for p := scope.From; p <= scope.To; p++ {
			ports = append(ports, p)
		}
		return ports, nil

	case PortsAll:
		var ports []int
		for p := 1; p <= 65535; p++ {
			ports = append(ports, p)
		}
		return ports, nil

	default:
		return nil, fmt.Errorf("unknown port scope")
	}
}

func ParsePortScope(token string) (PortScope, error) {
	token = strings.TrimSpace(strings.ToLower(token))

	////all
	if token == "all" {
		return PortScope{Type: PortsAll}, nil
	}

	// range: 1-1024
	if strings.Contains(token, "-") {
		parts := strings.Split(token, "-")
		if len(parts) != 2 {
			return PortScope{}, fmt.Errorf("invalid port range")
		}

		from, err := strconv.Atoi(parts[0])
		if err != nil {
			return PortScope{}, err
		}

		to, err := strconv.Atoi(parts[1])
		if err != nil {
			return PortScope{}, err
		}

		if from < 1 || to > 65535 || from > to {
			return PortScope{}, fmt.Errorf("invalid port range")
		}

		return PortScope{
			Type: PortsRange,
			From: from,
			To:   to,
		}, nil
	}

	// list: 22,80,443
	if strings.Contains(token, ",") {
		parts := strings.Split(token, ",")
		var ports []int

		for _, p := range parts {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil || port < 1 || port > 65535 {
				return PortScope{}, fmt.Errorf("invalid port: %s", p)
			}
			ports = append(ports, port)
		}

		return PortScope{
			Type:  PortsList,
			Ports: ports,
		}, nil
	}

	// single port
	port, err := strconv.Atoi(token)
	if err != nil || port < 1 || port > 65535 {
		return PortScope{}, fmt.Errorf("invalid port: %s", token)
	}

	return PortScope{
		Type:  PortSingle,
		Ports: []int{port},
	}, nil
}
